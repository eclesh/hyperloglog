package hyperloglog

import (
	"encoding/binary"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/DataDog/mmh3"
)

var buf32 = make([]byte, 4)
var buf64 = make([]byte, 8)
var buf128 = make([]byte, 16)

// Test that our abbreviated murmur hash works the same as upstream
func TestMurmur(t *testing.T) {
	for i := 0; i < 100; i++ {
		x := rand.Int31()
		binary.LittleEndian.PutUint32(buf32, uint32(x))
		hash := mmh3.Hash32(buf32)
		m := Murmur32(uint32(x))
		if hash != m {
			t.Errorf("Hash mismatch on 32 bit %d: expected 0x%X, got 0x%X\n", x, hash, m)
		}
	}

	for i := 0; i < 100; i++ {
		x := rand.Int63()
		binary.LittleEndian.PutUint64(buf64, uint64(x))
		hash := mmh3.Hash32(buf64)
		m := Murmur64(uint64(x))
		if hash != m {
			t.Errorf("Hash mismatch on 64 bit %d: expected 0x%X, got 0x%X\n", x, hash, m)
		}
	}

	for i := 0; i < 100; i++ {
		x := rand.Int63()
		y := rand.Int63()
		binary.LittleEndian.PutUint64(buf128, uint64(x))
		binary.LittleEndian.PutUint64(buf128[8:], uint64(y))
		hash := mmh3.Hash32(buf128)
		m := Murmur128(uint64(x), uint64(y))
		if hash != m {
			t.Errorf("Hash mismatch on 128 bit %d,%d: expected 0x%X, got 0x%X\n", x, y, hash, m)
		}
	}

	for i := 0; i < 100; i++ {
		key := randString((i % 15) + 5)
		hash := mmh3.Hash32([]byte(key))
		m := MurmurString(key)
		if hash != m {
			t.Errorf("Hash mismatch on key %s: expected 0x%X, got 0x%X\n", key, hash, m)
		}
	}
}

func randString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Benchmarks
func benchmarkMurmer64(b *testing.B, input []uint64) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			Murmur64(x)
		}
	}
}

func benchmarkMurmerString(b *testing.B, input []string) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			MurmurString(x)
		}
	}
}

func benchmarkHash32(b *testing.B, input []string) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, x := range input {
			b := *(*[]byte)(unsafe.Pointer(&x))
			mmh3.Hash32(b)
		}
	}
}

func Benchmark100Murmer64(b *testing.B) {
	input := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		input[i] = uint64(rand.Int63())
	}
	benchmarkMurmer64(b, input)
}

func Benchmark100MurmerString(b *testing.B) {
	input := make([]string, 100)
	for i := 0; i < 100; i++ {
		input[i] = randString((i % 15) + 5)
	}
	benchmarkMurmerString(b, input)
}

func Benchmark100Hash32(b *testing.B) {
	input := make([]string, 100)
	for i := 0; i < 100; i++ {
		input[i] = randString((i % 15) + 5)
	}
	benchmarkHash32(b, input)
}
