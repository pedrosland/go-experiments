package popcount_test

import (
	"testing"

	"github.com/pedrosland/go-experiments/popcount"
)

var table = []struct {
	binary   uint64
	expected int
}{
	{0x1, 1},
	{0x0, 0},
	{0x11111101, 7},
	{0xFFFFFFFFFFFFFFFF, 64},
}

func TestPopCount(t *testing.T) {
	runTest(popcount.PopCount, t)
}

func TestPopCountLoop(t *testing.T) {
	runTest(popcount.PopCountLoop, t)
}

func TestPopCountBitShift(t *testing.T) {
	runTest(popcount.PopCountBitShift, t)
}

func TestPopCountSubtract(t *testing.T) {
	runTest(popcount.PopCountSubtract, t)
}

func runTest(fn func(uint64) int, t *testing.T) {
	for i, test := range table {
		if result := fn(test.binary); result != test.expected {
			t.Errorf("test %d (%d) failed: expected %d got %d", i, test.binary, test.expected, result)
		}
	}
}

func BenchmarkPopCount(b *testing.B) {
	runBench(popcount.PopCount, b)
}

func BenchmarkPopCountLoop(b *testing.B) {
	runBench(popcount.PopCount, b)
}

func BenchmarkPopCountBitShift(b *testing.B) {
	runBench(popcount.PopCountBitShift, b)
}

func BenchmarkPopCountSubtract(b *testing.B) {
	runBench(popcount.PopCountSubtract, b)
}

func runBench(fn func(uint64) int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range table {
			// Assume results are correct
			fn(test.binary)
		}
	}
}
