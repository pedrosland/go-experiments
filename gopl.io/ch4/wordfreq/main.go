package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	words := make(map[string]int) // counts of Unicode chars in categories

	reader := bufio.NewScanner(os.Stdin)
	reader.Split(bufio.ScanWords)

	for reader.Scan() {
		words[reader.Text()]++
	}

	if reader.Err() != nil {
		fmt.Println("unable to read data", reader.Err())
		os.Exit(1)
	}

	fmt.Printf("word\t\tcount\n")
	for word, num := range words {
		fmt.Printf("%s\t\t%d\n", word, num)
	}
}
