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
	// only fields that are useful to us are unmarshalled
	ID         int `json:"num"`
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

	// Option A: read entries from disk at query time
	// runQuery(indexDir, searchQuery)

	// Option B: read entries into index and query index for results
	// Entries still have to be loaded from file initially and before displaying.
	index := prepareIndex(indexDir)
	runQueryOnIndex(index, searchQuery)

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
		comic, err := loadComic(fileInfo.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		// Note that search is case sensitive.
		if strings.Contains(comic.Title, query) || strings.Contains(comic.Transcript, query) {
			printMatch(comic)
			numMatches++
		}
	}
}

func prepareIndex(indexDir *os.File) *Index {
	indexFiles, err := indexDir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	index := NewIndex()

	for _, fileInfo := range indexFiles {
		filePath := path.Join(indexDirName, fileInfo.Name())
		comic, err := loadComic(filePath)
		if err != nil {
			log.Println(err)
			continue
		}

		index.addComic(*comic)
	}

	return index
}

func runQueryOnIndex(index *Index, query string) {
	words := strings.Split(query, " ")
	results := []int{}

	for i, word := range words {
		tmpResults := index.SearchWord(word)

		if len(tmpResults) == 0 {
			return
		}

		if i == 0 {
			results = tmpResults
			continue
		}

		oldResults := results
		results = []int{}

		for _, oldResult := range oldResults {
			for _, tmpResult := range tmpResults {
				if tmpResult == oldResult {
					results = append(results, oldResult)
					break
				}
			}
		}
	}

	for _, result := range results {
		filePath := path.Join(indexDirName, fmt.Sprintf("%d.info.0.json", result))
		comic, err := loadComic(filePath)
		if err != nil {
			log.Println(err)
			continue
		}

		printMatch(comic)
		numMatches++
	}
}

func printMatch(comic *Comic) {
	separator := "--------------\n"

	if numMatches == 0 {
		separator = ""
	}

	fmt.Printf(`%sTitle: %s
URL: %s
%s
`, separator, comic.Title, comic.ImageURL, comic.Transcript)
}

func loadComic(filePath string) (*Comic, error) {
	comic := Comic{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("warning: could not open file %s: %s\n", filePath, err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&comic)
	if err != nil {
		return nil, fmt.Errorf("warning: bad comic: invalid json in file %s: %s\n", filePath, err)
	}

	return &comic, nil
}
