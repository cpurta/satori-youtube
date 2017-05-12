package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/cpurta/satori/satori-youtube/cmd/internal/youtube"
	cache "github.com/patrickmn/go-cache"
)

type Crawler struct {
	publishChan chan json.RawMessage
	urlChan     chan string
	fetcher     Fetcher
	cache       *cache.Cache
	client      *youtube.VideoAPIClient
	shutdown    bool
}

func NewCrawler(pc chan json.RawMessage, urls chan string, fetcher Fetcher, c *cache.Cache, client *youtube.VideoAPIClient) *Crawler {
	return &Crawler{
		publishChan: pc,
		urlChan:     urls,
		fetcher:     fetcher,
		cache:       c,
		client:      client,
		shutdown:    false,
	}
}

// Crawl will pull from the url channel and attempt to pull data on the video if there is
// a parameter and then attempt to pull all url and only push those that conform to our
// video url regex
func (crawler *Crawler) Crawl() {
	for {
		select {
		case u := <-crawler.urlChan:
			cacheLock.Lock()
			_, crawled := crawler.cache.Get(u)
			cacheLock.Unlock()

			if !crawled {
				log.Println("Crawling", u)

				cacheLock.Lock()
				crawler.cache.Add(u, true, cache.DefaultExpiration)
				cacheLock.Unlock()

				urlquery, _ := url.Parse(u)
				videoResp, err := crawler.client.ListReqeust(urlquery.Query().Get("v"))
				if err != nil {
					log.Println("Error getting video snippet from YouTube API:", err.Error())
				}

				if err == nil && len(videoResp.Items) > 0 {
					message, _ := json.Marshal(videoResp.Items[0])
					pubChan <- message
				}

				urls, err := crawler.fetcher.Fetch(u)

				if err != nil {
					log.Printf("Error fetching results from %s: %s\n", u, err.Error())
				}

				for _, newURL := range urls {
					if !urlTest.MatchString(newURL) {
						newURL = "http://" + urlquery.Host + newURL
					}

					if validURL.MatchString(newURL) && urlTest.MatchString(newURL) {
						crawler.urlChan <- newURL
						log.Printf("Pushed %s to url channel\n", newURL)
					}
				}
			}
		default:
			if crawler.shutdown {
				return
			}
		}
	}
}

func (crawler *Crawler) Shutdown() {
	crawler.shutdown = true
}
