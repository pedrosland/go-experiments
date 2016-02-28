// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 133.

// Outline prints the outline of an HTML document tree.
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(os.Stdout, url)
	}
}

func outline(w io.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	//!+call
	forEachNode(w, doc, startElement, endElement)
	//!-call

	return nil
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(w io.Writer, n *html.Node, pre, post func(w io.Writer, n *html.Node)) {
	if pre != nil {
		pre(w, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(w, c, pre, post)
	}

	if post != nil {
		post(w, n)
	}
}

//!-forEachNode

//!+startend
var depth int
var inline bool

// Note: this changes the meaning of the HTML by inserting spaces
func startElement(w io.Writer, n *html.Node) {
	if n.Type == html.ElementNode {
		if n.FirstChild == nil {
			// Element has no children
			fmt.Fprintf(w, "%*s<%s%s/>\n", depth*2, "", n.Data, getAttrs(n))
			return
		}

		fmt.Fprintf(w, "%*s<%s%s>\n", depth*2, "", n.Data, getAttrs(n))
		depth++
	} else if n.Type == html.CommentNode {
		fmt.Fprintf(w, "%*s<!--%s-->\n", depth*2, "", n.Data)
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if len(text) == 0 {
			return
		}
		fmt.Fprintf(w, "%*s%s\n", depth*2, "", text)
	}
}

func endElement(w io.Writer, n *html.Node) {
	if n.Type == html.ElementNode {
		if n.FirstChild == nil {
			// Element has no children - already printed self-closing tag above
			return
		}

		depth--
		fmt.Fprintf(w, "%*s</%s>\n", depth*2, "", n.Data)
	}
}

func getAttrs(n *html.Node) string {
	if len(n.Attr) == 0 {
		return ""
	}

	var attrs []string
	for _, attr := range n.Attr {
		buff := bytes.Buffer{}
		if attr.Namespace != "" {
			buff.WriteString(attr.Namespace)
			buff.WriteRune(':')
		}

		buff.WriteString(attr.Key)
		buff.WriteString("=\"")
		buff.WriteString(html.EscapeString(attr.Val))
		buff.WriteRune('"')

		attrs = append(attrs, buff.String())
	}

	return fmt.Sprintf(" %s", strings.Join(attrs, " "))
}

//!-startend
