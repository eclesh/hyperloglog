package hyperloglog

// This file implements the murmur3 32-bit hash on 32bit and 64bit integers
// for little endian machines only with no heap allocation.  If you are using
// HLL to count integer IDs on intel machines, this is your huckleberry.

// Murmur32 implements a fast version of the murmur hash function for uint32 for
// little endian machines.  Suitable for adding 32bit integers to a HLL counter.
func Murmur32(i uint32) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32
	k = i
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second part
	h ^= 4
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

// Murmur64 implements a fast version of the murmur hash function for uint64 for
// little endian machines.  Suitable for adding 64bit integers to a HLL counter.
func Murmur64(i uint64) uint32 {
	var c1, c2 uint32 = 0xcc9e2d51, 0x1b873593
	var h, k uint32
	//first 4-byte chunk
	k = uint32(i)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second 4-byte chunk
	k = uint32(i >> 32)
	k *= c1
	k = (k << 15) | (k >> (32 - 15))
	k *= c2
	h ^= k
	h = (h << 13) | (h >> (32 - 13))
	h = (h * 5) + 0xe6546b64
	// second part
	h ^= 8
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}
