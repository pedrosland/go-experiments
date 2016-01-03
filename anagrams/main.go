package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Printf("  %t\n", anagram(os.Args[1], os.Args[2]))
}

func anagram(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, r := range s1 {
		if !strings.ContainsRune(s2, r) {
			return false
		}
	}

	return true
}
