// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 133.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	// HTML element ID, Web page URL
	run(os.Args[1], os.Args[2])
}

func run(id, url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	elem := ElementByID(doc, id)
	if elem == nil {
		fmt.Printf("could not find element with id: %s\n", id)
		return
	}

	fmt.Printf("found elem with id \"%s\": <%s>\n", id, elem.Data)

	return
}

func ElementByID(doc *html.Node, id string) *html.Node {
	fn := func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return false
			}
		}

		return true
	}

	return forEachNode(doc, fn, nil)
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) *html.Node {
	if pre != nil {
		cont := pre(n)
		if !cont {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := forEachNode(c, pre, post)
		if result != nil {
			return result
		}
	}

	if post != nil {
		cont := post(n)
		if !cont {
			return n
		}
	}

	return nil
}
