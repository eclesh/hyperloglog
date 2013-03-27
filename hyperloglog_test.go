package hyperloglog

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"math"
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

func testHyperLogLog(t *testing.T, n, low_b, high_b int) {
	words := dictionary(n)
	bad := 0
	n_words := uint64(len(words))
	for i := low_b; i < high_b; i++ {
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

		expected_error := 1.04 / math.Sqrt(float64(m))
		actual_error := math.Abs(geterror(n_words, h.Count()))

		if actual_error > expected_error {
			bad++
			t.Logf("m=%d: error=%.5f, expected <%.5f; actual=%d, estimated=%d\n",
				m, actual_error, expected_error, n_words, h.Count())
		}

	}
	t.Logf("%d of %d tests exceeded estimated error", bad, high_b-low_b)
}

func TestHyperLogLogSmall(t *testing.T) {
	testHyperLogLog(t, 5, 4, 17)
}

func TestHyperLogLogBig(t *testing.T) {
	testHyperLogLog(t, 0, 4, 17)
}
