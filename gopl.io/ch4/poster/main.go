package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	title := strings.Join(os.Args[1:], " ")

	movie, err := GetMovieInfo(title)
	if err != nil {
		fmt.Printf("error downloading movie info: %s\n", err)
		os.Exit(2)
	}

	posterReader, err := movie.GetPoster()
	if err != nil {
		fmt.Printf("error downloading poster: %s\n", err)
		os.Exit(2)
	}

	defer posterReader.Close()

	file, err := os.Create(strings.Title(title) + ".jpeg")
	if err != nil {
		fmt.Printf("unable to create poster file: %s\n", err)
		os.Exit(2)
	}

	defer file.Close()

	_, err = io.Copy(file, posterReader)
	if err != nil {
		fmt.Printf("error downloading file: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("Poster written to %s\n", file.Name())
}
