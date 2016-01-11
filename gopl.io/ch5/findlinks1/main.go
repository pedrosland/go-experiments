// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 122.
//!+main

// Findlinks1 prints the links in an HTML document read from standard input.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}

	// fmt.Println("Tag\tCount")
	// for tag, count := range countTags(make(map[string]int), doc) {
	// 	fmt.Printf("%s\t%d\n", tag, count)
	// }

	// printTextNodes(doc, "")
}

//!-main

//!+visit
// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		} else if n.Data == "link" {
			isStylesheet := false
			href := ""
			for _, attr := range n.Attr {
				if attr.Key == "rel" && attr.Val == "stylesheet" {
					isStylesheet = true
					if href != "" {
						links = append(links, href)
						break
					}
				} else if attr.Key == "href" {
					if isStylesheet {
						links = append(links, attr.Val)
						break
					} else {
						href = attr.Val
					}
				}
			}
		} else if n.Data == "script" || n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					links = append(links, attr.Val)
				}
			}
		}
	}
	if n.FirstChild != nil {
		links = visit(links, n.FirstChild)
	}
	if n.NextSibling != nil {
		links = visit(links, n.NextSibling)
	}
	return links
}

//!-visit

func countTags(tags map[string]int, n *html.Node) map[string]int {
	if n.Type == html.ElementNode {
		tags[n.Data]++
	}

	if n.FirstChild != nil {
		countTags(tags, n.FirstChild)
	}
	if n.NextSibling != nil {
		countTags(tags, n.NextSibling)
	}

	return tags
}

func printTextNodes(n *html.Node, tagName string) {
	if n.Type == html.TextNode {
		fmt.Printf("%s: %s\n", tagName, n.Data)
		return
	}

	if n.FirstChild != nil && n.Data != "script" && n.Data != "style" {
		printTextNodes(n.FirstChild, n.Data)
	}
	if n.NextSibling != nil {
		printTextNodes(n.NextSibling, n.Data)
	}
}

/*
//!+html
package html

type Node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func Parse(r io.Reader) (*Node, error)
//!-html
*/
