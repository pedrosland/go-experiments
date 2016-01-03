package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	diffCount := 0

	for i, b := range c1 {
		diff := uint8(b) ^ uint8(c2[i])

		diffCount += PopCountSubtract(diff)
	}

	if diffCount > 0 {
		fmt.Printf("%d different bits in hashes\n", diffCount)
	} else {
		fmt.Printf("The hashes are identical. 0 bits different\n")
	}
}

// PopCountSubtract returns the population count (number of set bits) of x.
// Not the fastest way but it is short
func PopCountSubtract(x uint8) (count int) {
	for x > 0 {
		x = x & (x - 1)
		count++
	}
	return
}
