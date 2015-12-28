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
	for i, test := range table {
		if result := popcount.PopCount(test.binary); result != test.expected {
			t.Errorf("test %d (%d) failed: expected %d got %d", i, test.binary, test.expected, result)
		}
	}
}

func TestPopCountLoop(t *testing.T) {
	for i, test := range table {
		if result := popcount.PopCountLoop(test.binary); result != test.expected {
			t.Errorf("test %d (%d) failed: expected %d got %d", i, test.binary, test.expected, result)
		}
	}
}

func TestPopCountBitShift(t *testing.T) {
	for i, test := range table {
		if result := popcount.PopCountBitShift(test.binary); result != test.expected {
			t.Errorf("test %d (%d) failed: expected %d got %d", i, test.binary, test.expected, result)
		}
	}
}

func TestPopCountSubtract(t *testing.T) {
	for i, test := range table {
		if result := popcount.PopCountSubtract(test.binary); result != test.expected {
			t.Errorf("test %d (%d) failed: expected %d got %d", i, test.binary, test.expected, result)
		}
	}
}
