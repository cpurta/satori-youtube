package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/url"
	"regexp"
	"sync"

	"github.com/cpurta/go-crawler/cmd/internal/config"
	"github.com/cpurta/go-crawler/cmd/internal/satori"
	"github.com/cpurta/go-crawler/cmd/internal/youtube"
	"github.com/willf/bloom"
)

var (
	seedUrl  string
	search   string
	depth    int
	crawlers int

	validURL *regexp.Regexp

	filterLock sync.Mutex

	throttle chan int
	pubChan  chan json.RawMessage

	crawlError = errors.New("already crawled")

	urlTest = regexp.MustCompile(`^((http[s]?):\/)?\/?([^:\/\s]+)((\/\w+)*\/)([\w\-\.]+[^#?\s]+)(.*)?(#[\w\-]+)?$`)
)

func main() {
	initFlags()
	if err := checkFlags(); err != nil {
		log.Printf("Error: %s", err.Error())
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Unable to load environment configuration:", err.Error())
	}

	throttle = make(chan int, crawlers)
	filterLock = sync.Mutex{}

	if search != "" {
		validURL = regexp.MustCompile(search)
	}

	filter := bloom.New(100000, 100)

	pubChan = make(chan json.RawMessage)

	publisher := satori.NewPublisher(config, pubChan)
	publisher.Start()
	go publisher.Publish()

	client := youtube.NewVideoAPIClient(config.YoutubeAuth)

	fetcher := URLFetcher{}
	log.Println("Starting crawl...")
	Crawl(seedUrl, depth, fetcher, filter, client)

	publisher.Shutdown()
}

func Crawl(searchUrl string, depth int, fetcher Fetcher, filter *bloom.BloomFilter, client *youtube.VideoAPIClient) {
	throttle <- 1

	if depth <= 0 {
		return
	}

	filterLock.Lock()
	crawled := filter.TestString(searchUrl)
	filterLock.Unlock()

	var wg sync.WaitGroup
	if !crawled {
		filterLock.Lock()
		filter.AddString(searchUrl)
		filterLock.Unlock()

		host, _ := url.Parse(searchUrl)

		videoResp, err := client.ListReqeust(host.Query().Get("v"))
		if err != nil {
			log.Println("Error getting video snippet from YouTube API:", err.Error())
		}

		if err == nil && len(videoResp.Items) > 0 {
			message, _ := json.Marshal(videoResp.Items[0])
			pubChan <- message
		}

		urls, err := fetcher.Fetch(searchUrl)
		if err != nil {
			log.Printf("Error fetching results from %s: %s\n", searchUrl, err.Error())
		}

		for _, u := range urls {
			if !urlTest.MatchString(u) {
				u = "http://" + host.Host + u
			}

			if validURL.MatchString(u) && urlTest.MatchString(u) {
				wg.Add(1)
				go func(u string, depth int, fetcher Fetcher, filter *bloom.BloomFilter) {
					defer wg.Done()
					Crawl(u, depth-1, fetcher, filter, client)
				}(u, depth, fetcher, filter)
			}
		}
	}

	<-throttle
	wg.Wait()
}

func initFlags() {
	flag.IntVar(&depth, "depth", 0, "The depth of how far the crawler will search in the network graph. Must be greater than 0.")
	flag.StringVar(&seedUrl, "seed-url", "", "The root url from which the crawler will look for network links.")
	flag.StringVar(&search, "search", "^.*$", `Regex that will be used against the urls crawled. Only urls matching the regex will be crawled. e.g. ^http(s)?://cnn.com\?+([0-9a-zA-Z]=[0-9a-zA-Z])$`)
	flag.IntVar(&crawlers, "crawlers", 10, "The number of concurrent crawling routines that will be used to crawl the web. Default: 10")
}

func checkFlags() error {
	flag.Parse()
	if seedUrl == "" {
		return errors.New("url flag cannot be empty")
	}
	if depth <= 0 {
		return errors.New("depth cannot be less than to equal to 0")
	}

	return nil
}
