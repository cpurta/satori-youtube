package youtube

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type VideoData struct {
	ID         string     `json:"id"`
	Snippet    Snippet    `json:"snippet"`
	Statistics Statistics `json:"statistics,omitempty"`
}

type Snippet struct {
	PublishedAt  string   `json:"published_at"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ChannelTitle string   `json:"channel_title"`
	Tags         []string `json:"tags"`
}

type Statistics struct {
	ViewCount    string `json:"view_count"`
	LikeCount    string `json:"like_count"`
	DislikeCount string `json:"dislike_count"`
	CommentCount string `json:"comment_count,omitempty"`
}

func ScrapeSnippetData(doc *goquery.Document) Snippet {
	published := doc.Find("strong.watch-time-text").Contents().Text()
	channelTitle := doc.Find("div.yt-user-info a").Contents().Text()
	title := doc.Find("span.watch-title").Contents().Text()
	description := doc.Find("p#eow-description").Contents().Text()

	// go and collect all the catgory tags
	tags := make([]string, 0)
	doc.Find("ul.watch-info-tag-list").Each(func(i int, item *goquery.Selection) {
		item.Find("li").Each(func(i int, item *goquery.Selection) {
			tag := item.Find("a").Contents().Text()
			if tag != "" {
				tags = append(tags, tag)
			}
		})
	})

	published = sanitizeString(published, []string{"Published on "}, "")

	snippet := Snippet{
		PublishedAt:  published,
		ChannelTitle: channelTitle,
		Title:        title,
		Description:  description,
		Tags:         tags,
	}

	return snippet
}

func ScrapeStatisticsData(doc *goquery.Document) Statistics {
	views := doc.Find(".watch-view-count").Contents().Text()
	likes := doc.Find("button.like-button-renderer-like-button span.yt-uix-button-content").Contents().Text()
	dislikes := doc.Find("button.like-button-renderer-dislike-button span.yt-uix-button-content").Contents().Text()
	// commentCount := doc.Find("h2.comment-section-header-renderer").Text()

	views = sanitizeString(views, []string{" views", ","}, "")
	likes = sanitizeString(likes, []string{","}, "")
	dislikes = sanitizeString(dislikes, []string{","}, "")

	stats := Statistics{
		ViewCount:    views,
		LikeCount:    likes[:len(likes)/2], // for some reason scraping like and dislikes concats the LikeCount+1 to the count
		DislikeCount: dislikes[:len(dislikes)/2],
		// CommentCount: commentCount,
	}

	return stats
}

func sanitizeString(target string, sanitize []string, newStr string) string {
	sanitizedStr := target
	for _, str := range sanitize {
		sanitizedStr = strings.Replace(sanitizedStr, str, newStr, -1)
	}

	return sanitizedStr
}
