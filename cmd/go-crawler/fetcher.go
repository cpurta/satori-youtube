package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
)

type Fetcher interface {
	Fetch(url string) (urls []string, err error)
}

type URLResult struct {
	URLs []string
	Body string
}

type URLFetcher map[string]*URLResult

func (f URLFetcher) Fetch(url string) (urls []string, err error) {
	if url == "" {
		return nil, errors.New("Empty URL provided")
	}

	// grab the response from the site...
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		tokenizer := html.NewTokenizer(bytes.NewReader(body))

		for {
			tt := tokenizer.Next()
			switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken:
				token := tokenizer.Token()

				// go and gather any links that we find on the page
				anchor := token.Data == "a"
				// If not an anchor tag then move on
				if !anchor {
					continue
				} else {
					ref := getHref(token)
					urls = append(urls, ref)
				}
			} // end switch
		} // end for
	} // end if

	return urls, nil
}

func getHref(token html.Token) string {
	for _, a := range token.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}

	return ""
}
