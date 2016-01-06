// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses a slice.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	//!+array
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverseArr(&a)
	fmt.Println(a) // "[5 4 3 2 1 0]"
	reverse(a[:])
	fmt.Println(a) // "[0 1 2 3 4 5]"
	//!-array

	//!+slice
	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	reverse(s[:2])
	reverse(s[2:])
	reverse(s)
	fmt.Println(s) // "[2 3 4 5 0 1]"
	//!-slice

	s = []int{0, 1, 2, 3, 4, 5}
	s = rotate(s, 2)
	fmt.Println(s)

	strS := []string{"a", "a", "b", "c", "c", "c"}
	fmt.Println(strS)
	strS = removeDuplicates(strS)
	fmt.Println(strS)

	str := "a b　c　　d"
	fmt.Println(str)
	b := []byte(str)
	fmt.Println(b)
	b = replaceDuplicateSpaces(b)
	fmt.Println(b)
	fmt.Println(string(b))

	// Interactive test of reverse.
	input := bufio.NewScanner(os.Stdin)
outer:
	for input.Scan() {
		var ints []int
		for _, s := range strings.Fields(input.Text()) {
			x, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue outer
			}
			ints = append(ints, int(x))
		}
		reverse(ints)
		fmt.Printf("%v\n", ints)
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!+rev
// reverse reverses a slice of ints in place.
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

//!-rev

//
// -------------------------------------------------------
//

// reverseArr reverses an array of 6 ints
func reverseArr(a *[6]int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func rotate(s []int, count int) []int {
	numItems := len(s)
	items := make([]int, numItems, numItems)
	newPos := numItems - count

	copy(items, s[count:])
	copy(items[newPos:], s[:count])

	return items
}

// removeDuplicates removes adjacent duplicate strings in a slice
// This modifies the original slices's array
func removeDuplicates(s []string) []string {
	var lastStr string
	i := 0

	for _, str := range s {
		if str != lastStr {
			s[i] = str
			i++
			lastStr = str
		}
	}

	return s[:i]
}

func replaceDuplicateSpaces(b []byte) []byte {
	// convert to []string so that `range` iterates over chars and not bytes
	str := string(b)
	lastCharSpace := false
	i := 0

	for _, char := range str {
		if unicode.IsSpace(char) {
			if !lastCharSpace {
				b[i] = byte(' ') // ASCII space
				lastCharSpace = true
				i++
			}
			continue
		} else if lastCharSpace {
			lastCharSpace = false
		}

		charStr := string(char)
		charLen := len(charStr)
		copy(b[i:], charStr)
		i += charLen
	}

	return b[:i]
}
