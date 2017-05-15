package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/cpurta/satori/satori-youtube/cmd/internal/config"
	"github.com/cpurta/satori/satori-youtube/cmd/internal/satori"
	"github.com/cpurta/satori/satori-youtube/cmd/internal/youtube"
	"github.com/patrickmn/go-cache"
)

var (
	seedUrl  string
	search   string
	depth    int
	routines int

	validURL *regexp.Regexp

	cacheLock sync.Mutex

	pubChan chan json.RawMessage

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

	cacheLock = sync.Mutex{}

	if search != "" {
		validURL = regexp.MustCompile(search)
	}

	cache := cache.New(5*time.Minute, 5*time.Minute)

	pubChan = make(chan json.RawMessage)
	urls := make(chan string, 250000)

	publisher := satori.NewPublisher(config, pubChan)
	publisher.Start()
	go publisher.Publish()

	client := youtube.NewVideoAPIClient(config.YoutubeAuth)

	go cleanURLs(urls)

	fetcher := URLFetcher{}
	log.Println("Starting crawl...")

	crawlers := make([]*Crawler, routines)

	var wg sync.WaitGroup
	for i := 0; i < routines; i++ {
		crawlers[i] = NewCrawler(pubChan, urls, fetcher, cache, client)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			crawlers[i].Crawl()
		}(i)
	}

	urls <- seedUrl

	wg.Wait()

	publisher.Shutdown()

	close(pubChan)
}

func cleanURLs(urls chan string) {
	for {
		if len(urls) == cap(urls) {
			log.Println("Dumping half of the urls from the channel to make room")
			for i := 0; i < cap(urls)/2; i++ {
				<-urls
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func initFlags() {
	flag.IntVar(&depth, "depth", 0, "The depth of how far the crawler will search in the network graph. Must be greater than 0.")
	flag.StringVar(&seedUrl, "seed-url", "", "The root url from which the crawler will look for network links.")
	flag.StringVar(&search, "search", "^.*$", `Regex that will be used against the urls crawled. Only urls matching the regex will be crawled. e.g. ^http(s)?://cnn.com\?+([0-9a-zA-Z]=[0-9a-zA-Z])$`)
	flag.IntVar(&routines, "crawlers", 10, "The number of concurrent crawling routines that will be used to crawl the web. Default: 10")
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
