package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
)

func main() {
	var digest hash.Hash

	hashType := flag.String("hash", "", "Which hashing algorithm to use. Either 256 (default), 384 or 512.")
	flag.Parse()

	switch *hashType {
	case "":
		fallthrough
	case "256":
		digest = sha256.New()
	case "384":
		digest = sha512.New384()
	case "512":
		digest = sha512.New()
	default:
		flag.Usage()
		os.Exit(64) // Invalid usage
	}

	io.Copy(digest, os.Stdin)

	hash := digest.Sum([]byte{})

	fmt.Printf("%x\n", hash)
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
