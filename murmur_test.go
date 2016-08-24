package hyperloglog

import (
	"encoding/binary"
	"math/rand"
	"testing"

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
}
