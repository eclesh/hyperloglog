// Package hyperloglog implements the HyperLogLog algorithm for
// cardinality estimation. In English: it counts things. It counts
// things using very small amounts of memory compared to the number of
// objects it is counting.
//
// For a full description of the algorithm, see the paper HyperLogLog:
// the analysis of a near-optimal cardinality estimation algorithm by
// Flajolet, et. al.
package hyperloglog

import (
	"fmt"
	"math"
)

var (
	exp32 = math.Pow(2, 32)
)

type HyperLogLog struct {
	m         uint    // Number of registers
	b         uint32  // Number of bits used to determine register index
	alpha     float64 // Bias correction constant
	registers []uint8
}

// Compute bias correction alpha_m.
func get_alpha(m uint) (result float64) {
	switch m {
	case 16:
		result = 0.673
	case 32:
		result = 0.697
	case 64:
		result = 0.709
	default:
		result = 0.7213 / (1.0 + 1.079/float64(m))
	}
	return result
}

// Return a new HyperLogLog with the given number of registers. More
// registers leads to lower error in your estimated count, at the
// expense of memory.
//
// Choose a power of two number of registers, depending on the amount
// of memory you're willing to use and the error you're willing to
// tolerate. Each register uses one byte of memory.
//
// Approximate error will be:
//     1.04 / sqrt(registers)
//
func New(registers uint) (*HyperLogLog, error) {
	if (registers & (registers - 1)) != 0 {
		return nil, fmt.Errorf("number of registers %d not a power of two", registers)
	}
	h := &HyperLogLog{}
	h.m = registers
	h.b = uint32(math.Ceil(math.Log2(float64(registers))))
	h.alpha = get_alpha(registers)
	h.Reset()
	return h, nil
}

// Reset all internal variables and set the count to zero.
func (h *HyperLogLog) Reset() {
	h.registers = make([]uint8, h.m)
}

// Calculate the position of the leftmost 1-bit.
func rho(val uint32, max uint32) uint8 {
	r := uint32(1)
	for val&0x80000000 == 0 && r <= max {
		r++
		val <<= 1
	}
	return uint8(r)
}

// Add to the count. val should be a 32 bit unsigned integer from a
// good hash function.
func (h *HyperLogLog) Add(val uint32) {
	k := 32 - h.b
	r := rho(val<<h.b, k)
	j := val >> uint(k)
	if r > h.registers[j] {
		h.registers[j] = r
	}
}

// Get the estimated count.
func (h *HyperLogLog) Count() uint64 {
	sum := 0.0
	for _, val := range h.registers {
		sum += 1.0 / math.Pow(2.0, float64(val))
	}
	estimate := h.alpha * float64(h.m*h.m) / sum
	if estimate <= 5.0/2.0*float64(h.m) {
		// Small range correction
		v := 0
		for _, r := range h.registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = float64(h.m) * math.Log(float64(h.m)/float64(v))
		}
	} else if estimate > 1.0/30.0*exp32 {
		// Large range correction
		estimate = -exp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge another HyperLogLog into this one. The number of registers in
// each must be the same.
func (h1 *HyperLogLog) Merge(h2 *HyperLogLog) error {
	if h1.m != h2.m {
		return fmt.Errorf("number of registers doesn't match: %d != %d",
			h1.m, h2.m)
	}
	for j, r := range h2.registers {
		if r > h1.registers[j] {
			h1.registers[j] = r
		}
	}
	return nil
}
