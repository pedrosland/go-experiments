// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 112.
//!+

// Issues prints a table of GitHub issues matching the search terms.
//
// Usage: ./issues [github search terms...]
// eg. ./issues repo:golang/go is:open json encoder
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

//!+
func main() {
	var binLt1Month, binLt1Year, binGt1Year []*github.Issue

	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues\n", result.TotalCount)

	now := time.Now()
	newMonth := now.Month() - 1
	if newMonth <= 0 {
		newMonth = time.December
	}
	oneMonthAgo := time.Date(now.Year(), now.Month()-1, now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	oneYearAgo := time.Date(now.Year()-1, now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	for _, item := range result.Items {
		if t := item.CreatedAt; t.After(oneMonthAgo) {
			binLt1Month = append(binLt1Month, item)
		} else if t.After(oneYearAgo) {
			binLt1Year = append(binLt1Year, item)
		} else {
			binGt1Year = append(binGt1Year, item)
		}
	}

	printResults(binLt1Month, "Less than 1 month old")
	printResults(binLt1Year, "Less than 1 year old")
	printResults(binGt1Year, "Over 1 year old")
}

func printResults(s []*github.Issue, heading string) {
	if len(s) > 0 {
		fmt.Printf("\n%s:\n", heading)
		for _, item := range s {
			fmt.Printf("#%-5d %9.9s %.55s\n",
				item.Number, item.User.Login, item.Title)
		}
	}
}

//!-

/*
//!+textoutput
$ go build gopl.io/ch4/issues
$ ./issues repo:golang/go is:open json decoder
13 issues:
#5680    eaigner encoding/json: set key converter on en/decoder
#6050  gopherbot encoding/json: provide tokenizer
#8658  gopherbot encoding/json: use bufio
#8462  kortschak encoding/json: UnmarshalText confuses json.Unmarshal
#5901        rsc encoding/json: allow override type marshaling
#9812  klauspost encoding/json: string tag not symmetric
#7872  extempora encoding/json: Encoder internally buffers full output
#9650    cespare encoding/json: Decoding gives errPhase when unmarshalin
#6716  gopherbot encoding/json: include field name in unmarshal error me
#6901  lukescott encoding/json, encoding/xml: option to treat unknown fi
#6384    joeshaw encoding/json: encode precise floating point integers u
#6647    btracey x/tools/cmd/godoc: display type kind of each named type
#4237  gjemiller encoding/base64: URLEncoding padding is optional
//!-textoutput
*/
