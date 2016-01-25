package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	printer := make(chan string)
	go print(printer)

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		line := scan.Text()

		printer <- line
	}

	if err := scan.Err(); err != nil {
		panic(err)
	}
}

func print(input chan string) {
	for {
		msg := <-input
		fmt.Println(msg)
		// do work with msg
	}
}
