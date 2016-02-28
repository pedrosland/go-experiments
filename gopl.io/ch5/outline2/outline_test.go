package main

import (
	"bytes"
	"fmt"
	"testing"

	"golang.org/x/net/html"
)

// Test requires the internet
func TestGoplIo(t *testing.T) {
	b := &bytes.Buffer{}
	fmt.Fprint(b, "hello")
	err := outline(b, "http://gopl.io")
	if err != nil {
		t.Fatalf("producing HTML outline: %s", err)
	}

	// Normally the parser just returns a very basic html document instead of
	// returning an error so this doesn't mean much.
	_, err = html.Parse(b)
	if err != nil {
		t.Fatalf("invalid HTML: %s", err)
	}
}
