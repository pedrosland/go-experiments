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
	// discard spaces
	s1 = strings.Replace(s1, " ", "", -1)
	s2 = strings.Replace(s2, " ", "", -1)

	if len(s1) != len(s2) {
		return false
	}

	s1Chars := make(map[rune]int)

	for _, r := range s1 {
		// works because empty type is 0
		s1Chars[r]++
	}

	for r, count := range s1Chars {
		if strings.Count(s2, string(r)) != count {
			return false
		}
	}

	return true
}
