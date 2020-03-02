package hyperloglog

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"math/rand"
	"os"
	"testing"
)

// Return a dictionary up to n words. If n is zero, return the entire
// dictionary.
func dictionary(n int) []string {
	var words []string
	dict := "/usr/share/dict/words"
	f, err := os.Open(dict)
	if err != nil {
		fmt.Printf("can't open dictionary file '%s': %v\n", dict, err)
		os.Exit(1)
	}
	count := 0
	buf := bufio.NewReader(f)
	for {
		if n != 0 && count >= n {
			break
		}
		word, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		words = append(words, word)
		count++
	}
	f.Close()
	return words
}

func geterror(actual uint64, estimate uint64) (result float64) {
	return (float64(estimate) - float64(actual)) / float64(actual)
}

func testHyperLogLog(t *testing.T, n, lowB, highB int) {
	words := dictionary(n)
	bad := 0
	nWords := uint64(len(words))
	for i := lowB; i < highB; i++ {
		m := uint(math.Pow(2, float64(i)))

		h, err := New(m)
		if err != nil {
			t.Fatalf("can't make New(%d): %v", m, err)
		}

		hash := fnv.New32()
		for _, word := range words {
			hash.Write([]byte(word))
			h.Add(hash.Sum32())
			hash.Reset()
		}

		expectedError := 1.04 / math.Sqrt(float64(m))
		actualError := math.Abs(geterror(nWords, h.Count()))

		if actualError > expectedError {
			bad++
			t.Logf("m=%d: error=%.5f, expected <%.5f; actual=%d, estimated=%d\n",
				m, actualError, expectedError, nWords, h.Count())
		}

	}
	t.Logf("%d of %d tests exceeded estimated error", bad, highB-lowB)
}

func TestHyperLogLogSmall(t *testing.T) {
	testHyperLogLog(t, 5, 4, 17)
}

func TestHyperLogLogBig(t *testing.T) {
	testHyperLogLog(t, 0, 4, 17)
}

func testReset(t *testing.T, m uint, numObjects, runs int) {
	rand.Seed(101)

	h, err := New(m)
	if err != nil {
		t.Fatalf("can't make New(%d): %v", m, err)
	}

	for i := 0; i < runs; i++ {
		for j := 0; j < numObjects; j++ {
			h.Add(rand.Uint32())
		}

		oldRegisters := &h.Registers
		h.Reset()
		if oldRegisters != &h.Registers {
			t.Error("registers were reallocated")
		}
		for _, r := range h.Registers {
			if r != 0 {
				t.Error("register is not zeroed out after reset")
			}
		}
	}
}

func TestReset(t *testing.T) {
	testReset(t, 512, 1_000_000, 10)
}

func BenchmarkReset(b *testing.B) {
	m := uint(256)
	numObjects := 1000

	h, err := New(m)
	if err != nil {
		b.Fatalf("can't make New(%d): %v", m, err)
	}

	for n := 0; n < b.N; n++ {
		for i := 0; i < numObjects; i++ {
			h.Add(uint32(i))
		}
		h.Reset()
	}
}

func benchmarkCount(b *testing.B, registers int) {
	words := dictionary(0)
	m := uint(math.Pow(2, float64(registers)))

	h, err := New(m)
	if err != nil {
		return
	}

	hash := fnv.New32()
	for _, word := range words {
		hash.Write([]byte(word))
		h.Add(hash.Sum32())
		hash.Reset()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		h.Count()
	}
}

func BenchmarkCount4(b *testing.B) {
	benchmarkCount(b, 4)
}

func BenchmarkCount5(b *testing.B) {
	benchmarkCount(b, 5)
}

func BenchmarkCount6(b *testing.B) {
	benchmarkCount(b, 6)
}

func BenchmarkCount7(b *testing.B) {
	benchmarkCount(b, 7)
}

func BenchmarkCount8(b *testing.B) {
	benchmarkCount(b, 8)
}

func BenchmarkCount9(b *testing.B) {
	benchmarkCount(b, 9)
}

func BenchmarkCount10(b *testing.B) {
	benchmarkCount(b, 10)
}
