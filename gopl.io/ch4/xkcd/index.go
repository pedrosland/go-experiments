package main

import (
	"bytes"
	"strings"
	"unicode"
)

type Index struct {
	words map[string][]int
}

func NewIndex() *Index {
	return &Index{
		words: make(map[string][]int),
	}
}

func (idx *Index) addComic(comic Comic) {
	text := comic.Title + " " + comic.Transcript
	buf := new(bytes.Buffer)
	s := struct{}{}
	uniqueWords := make(map[string]struct{})

	// remove symbols
	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			buf.WriteString(strings.ToLower(string(char)))
		}

		if unicode.IsSpace(char) && buf.Len() > 0 {
			uniqueWords[buf.String()] = s
			buf.Reset()
		}
	}

	for word := range uniqueWords {
		idx.addWord(word, comic.ID)
	}
}

func (idx *Index) addWord(word string, comicID int) {
	comicIDs := idx.words[word]
	comicIDs = append(comicIDs, comicID)
	idx.words[word] = comicIDs
}

func (idx *Index) SearchWord(word string) []int {
	return idx.words[word]
}
