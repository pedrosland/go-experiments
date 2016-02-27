package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const APIURL = "https://www.omdbapi.com/"

type MovieInfo struct {
	Title  string
	Poster string
}

func GetMovieInfo(title string) (*MovieInfo, error) {
	values := make(url.Values)

	values.Set("r", "json")
	values.Set("t", title)

	resp, err := http.Get(fmt.Sprintf("%s?%s", APIURL, values.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	movie := &MovieInfo{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (movie *MovieInfo) GetPoster() (io.ReadCloser, error) {
	fmt.Println(movie.Poster)
	resp, err := http.Get(movie.Poster)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
