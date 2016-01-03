// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 73.

// Comma prints its argument numbers with a comma at each power of 1000.
//
// Example:
// 	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
// 	1
// 	12
// 	123
// 	1,234
// 	1,234,567,890
//
package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", nonRecursiveComma(os.Args[i]))
	}
}

//!+
// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func nonRecursiveComma(s string) string {
	var prefix byte
	buf := bytes.Buffer{}

	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		prefix = s[0]
		s = s[1:]
	}

	length := len(s)

	if length <= 3 {
		return string(prefix) + s
	}

	buf.WriteByte(prefix)

	for i := 0; i < length; i++ {
		if (length-i)%3 == 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(s[i : i+1])
	}

	return buf.String()
}

//!-