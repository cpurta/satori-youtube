package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type VideoAPIClient struct {
	AuthorizationKey string
}

type VideoListResponse struct {
	Kind     string             `json:"kind"`
	Etag     string             `json:"etag"`
	PageInfo map[string]float64 `json:"pageInfo"`
	Items    []Item             `json:"items"`
}

type Item struct {
	Kind       string     `json:"kind"`
	Etag       string     `json:"etag"`
	ID         string     `json:"id"`
	Snippet    Snippet    `json:"snippet"`
	Statistics Statistics `json:"statistics,omitempty"`
}

type Snippet struct {
	PublishedAt  string   `json:"publishedAt"`
	ChannelID    string   `json:"channelId"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ChannelTitle string   `json:"channelTitle"`
	Tags         []string `json:"tags"`
	CategoryID   string   `json:"categoryId"`
}

type Statistics struct {
	ViewCount     string `json:"viewCount"`
	LikeCount     string `json:"likeCount"`
	DislikeCount  string `json:"dislikeCount"`
	FavoriteCount string `json:"favoriteCount"`
	CommentCount  string `json:"commentCount"`
}

func NewVideoAPIClient(authKey string) *VideoAPIClient {
	return &VideoAPIClient{AuthorizationKey: authKey}
}

func (client *VideoAPIClient) ListReqeust(id string) (*VideoListResponse, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=%s&key=%s", id, client.AuthorizationKey)

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading YouTube API response body:", err.Error())
		return nil, err
	}

	var videoResponse VideoListResponse
	err = json.Unmarshal(body, &videoResponse)

	return &videoResponse, err
}

func ScrapeSnippetData(url string) (Snippet, error) {
	var snippet Snippet

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return snippet, err
	}

	published := doc.Find("strong.watch-time-text").Contents().Text()
	channelTitle := doc.Find("div.yt-user-info a").Contents().Text()
	title := doc.Find("span.watch-title").Contents().Text()
	description := doc.Find("p#eow-description").Contents().Text()

	snippet = Snippet{
		PublishedAt:  published,
		ChannelTitle: channelTitle,
		Title:        title,
		Description:  description,
	}

	return snippet, nil
}

func ScrapeStatisticsData(url string) (Statistics, error) {
	var stats Statistics

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return stats, err
	}

	viewCount := doc.Find(".watch-view-count").Contents().Text()
	likes := doc.Find("button.like-button-renderer-like-button span.yt-uix-button-content").Contents().Text()
	dislikes := doc.Find("button.like-button-renderer-dislike-button span.yt-uix-button-content").Contents().Text()
	commentCount := doc.Find("h2.comment-section-header-renderer b").Contents().Text()
	stats = Statistics{
		ViewCount:    viewCount,
		LikeCount:    likes,
		DislikeCount: dislikes,
		CommentCount: commentCount,
	}

	return stats, err
}
