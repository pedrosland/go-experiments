package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Comic struct {
	// only fields that are useful to us are scanned
	Transcript string
	Title      string
	ImageURL   string `json:"img"` // this isn't the CURL of the comic but the URL of the image
}

const indexDirName = "index"

var numMatches = 0

func main() {
	indexDir, err := os.Open(indexDirName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(`index does not exist. Try the following commands:
mkdir index
cd index
for i in $(seq 500 520); do; wget -O $i.info.0.json https://xkcd.com/$i/info.0.json; sleep 1; done
`)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}

	searchQuery := strings.Join(os.Args[1:], " ")

	runQuery(indexDir, searchQuery)

	if numMatches == 0 {
		log.Printf("no matches found for query \"%s\"\n", searchQuery)
	}
}

func runQuery(indexDir *os.File, query string) {
	indexFiles, err := indexDir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range indexFiles {
		comic := Comic{}

		file, err := os.Open(path.Join(indexDirName, fileInfo.Name()))
		if err != nil {
			log.Printf("warning: could not open file %s: %s\n", fileInfo.Name(), err)
			continue
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&comic)
		if err != nil {
			log.Printf("warning: bad comic: invalid json in file %s: %s\n", file.Name(), err)
			continue
		}

		// Note that search is case sensitive.
		if strings.Contains(comic.Title, query) || strings.Contains(comic.Transcript, query) {
			printMatch(comic)
			numMatches++
		}
	}
}

func printMatch(comic Comic) {
	separator := "--------------\n"

	if numMatches == 0 {
		separator = ""
	}

	fmt.Printf(`%sTitle: %s
URL: %s
%s
`, separator, comic.Title, comic.ImageURL, comic.Transcript)
}
