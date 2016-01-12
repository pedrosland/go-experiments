package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	words, images, err := CountWordsAndImages(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Words: %d\n", words)
	fmt.Printf("Images: %d\n", images)
}

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n.Type == html.TextNode {
		words += len(strings.Split(n.Data, " "))
		return
	}

	if n.Data == "img" {
		images++
	} else {
		if n.FirstChild != nil {
			w, i := countWordsAndImages(n.FirstChild)
			words += w
			images += i
		}
		if n.NextSibling != nil {
			w, i := countWordsAndImages(n.NextSibling)
			words += w
			images += i
		}
	}
	return
}
