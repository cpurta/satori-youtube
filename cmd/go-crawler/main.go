package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"runtime"
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
	crawlers int

	validURL *regexp.Regexp

	cacheLock sync.Mutex

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
	cacheLock = sync.Mutex{}

	if search != "" {
		validURL = regexp.MustCompile(search)
	}

	cache := cache.New(5*time.Minute, 5*time.Minute)

	pubChan = make(chan json.RawMessage)

	publisher := satori.NewPublisher(config, pubChan)
	publisher.Start()
	go publisher.Publish()

	client := youtube.NewVideoAPIClient(config.YoutubeAuth)

	fetcher := URLFetcher{}
	log.Println("Starting crawl...")

	go printStats()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	Crawl(seedUrl, depth, fetcher, cache, client, wg)

	wg.Wait()

	close(throttle)

	publisher.Shutdown()

	close(pubChan)
}

func Crawl(searchUrl string, depth int, fetcher Fetcher, c *cache.Cache, client *youtube.VideoAPIClient, wg *sync.WaitGroup) {
	throttle <- 1
	defer wg.Done()

	if depth <= 0 {
		return
	}

	cacheLock.Lock()
	_, crawled := c.Get(searchUrl)
	cacheLock.Unlock()

	if !crawled {
		cacheLock.Lock()
		c.Add(searchUrl, true, cache.DefaultExpiration)
		cacheLock.Unlock()

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
				<-throttle
				wg.Add(1)
				go Crawl(u, depth-1, fetcher, c, client, wg)
			}
		}
	}
}

func printStats() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Printf("Go Routines: %d\nHeap Allocated: %d\nTotal Memory Allocated: %d\n", runtime.NumGoroutine(), int64(m.HeapAlloc), int64(m.TotalAlloc))

		time.Sleep(time.Second * 10)
	}
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
