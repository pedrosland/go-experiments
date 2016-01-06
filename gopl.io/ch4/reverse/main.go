// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 86.

// Rev reverses a slice.
package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	//!+array
	a := [...]int{0, 1, 2, 3, 4, 5}
	fmt.Println("reverseArr()")
	reverseArr(&a)
	fmt.Println(a) // "[5 4 3 2 1 0]"
	fmt.Println("reverse")
	reverse(a[:])
	fmt.Println(a) // "[0 1 2 3 4 5]"
	//!-array
	fmt.Println()

	fmt.Println("reverse() x3 (rotate)")
	//!+slice
	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	reverse(s[:2])
	reverse(s[2:])
	reverse(s)
	fmt.Println(s) // "[2 3 4 5 0 1]"
	//!-slice
	fmt.Println()

	fmt.Println("rotate()")
	s = []int{0, 1, 2, 3, 4, 5}
	s = rotate(s, 2)
	fmt.Println(s)
	fmt.Println()

	fmt.Println("removeDuplicates()")
	strS := []string{"a", "a", "b", "c", "c", "c"}
	fmt.Println(strS)
	strS = removeDuplicates(strS)
	fmt.Println(strS)
	fmt.Println()

	fmt.Println("replaceDuplicateSpaces()")
	str := "a b　c　　d"
	fmt.Println(str)
	b := []byte(str)
	fmt.Println(b)
	b = replaceDuplicateSpaces(b)
	fmt.Println(b)
	fmt.Println(string(b))
	fmt.Println()

	fmt.Println("reverseByteString()")
	str = "a b　c　d　"
	fmt.Println(str)
	b = []byte(str)
	fmt.Println(b)
	reverseByteString(b)
	fmt.Println(b)
	fmt.Println(string(b))

	// Interactive test of reverse.
	// 	input := bufio.NewScanner(os.Stdin)
	// outer:
	// 	for input.Scan() {
	// 		var ints []int
	// 		for _, s := range strings.Fields(input.Text()) {
	// 			x, err := strconv.ParseInt(s, 10, 64)
	// 			if err != nil {
	// 				fmt.Fprintln(os.Stderr, err)
	// 				continue outer
	// 			}
	// 			ints = append(ints, int(x))
	// 		}
	// 		reverse(ints)
	// 		fmt.Printf("%v\n", ints)
	// 	}
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

// reverseByteString reverses a []byte that contains utf-8 encoded characters.
// Does operation in place while keeping memory usage to a minimum.
func reverseByteString(s []byte) {
	var firstRune, lastRune rune
	var firstSize, lastSize int

	length := len(s)

	for i, j := 0, length-1; i < j; i, j = i+lastSize, j-firstSize {
		firstRune, firstSize = utf8.DecodeRune(s[i:])
		lastRune, lastSize = utf8.DecodeLastRune(s[:j+1])

		// shift bytes to make the right space for lastRune
		copy(s[i+lastSize:], s[i+firstSize:j-lastSize+1])

		// copy lastRune to first
		copy(s[i:], []byte(string(lastRune)))

		// copy firstRune to last
		copy(s[j-firstSize+1:], []byte(string(firstRune)))
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

// replaceDuplicateSpaces squashes adjacent Unicode spaces into a single ASCII space
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
