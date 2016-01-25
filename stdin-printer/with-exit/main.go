package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	printer := make(chan string)
	reader := make(chan string)
	go print(printer)
	go read(reader)

loop:
	for {
		select {
		case line := <-reader:
			printer <- line
		case <-time.After(5 * time.Second):
			break loop
		}
	}
}

func print(input chan string) {
	for {
		msg := <-input
		fmt.Println(msg)
		// do work with msg
	}
}

func read(output chan string) {
	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		line := scan.Text()

		output <- line
	}

	if err := scan.Err(); err != nil {
		panic(err)
	}
}
