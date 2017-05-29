package youtube

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestScrapeSnippetData(t *testing.T) {
	doc, err := goquery.NewDocument("https://www.youtube.com/watch?v=i8YRtIHRIv0")
	if err != nil {
		t.Error(err)
	}

	snippet := ScrapeSnippetData(doc)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Title:", snippet.Title)
	fmt.Println("Description:", snippet.Description)
	fmt.Println("Channel Title:", snippet.ChannelTitle)
	fmt.Println("Published at:", snippet.PublishedAt)
	fmt.Println("Tags:", snippet.Tags)
}

func TestScrapeStatisticsData(t *testing.T) {
	doc, err := goquery.NewDocument("https://www.youtube.com/watch?v=i8YRtIHRIv0")
	if err != nil {
		t.Error(err)
	}

	stats := ScrapeStatisticsData(doc)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("View count:", stats.ViewCount)
	fmt.Println("Like count", stats.LikeCount)
	fmt.Println("Dislike count", stats.DislikeCount)
	fmt.Println("Comment count", stats.CommentCount)
}

func TestSanitizeString(t *testing.T) {
	viewCount := "15,309 views"

	santized := sanitizeString(viewCount, []string{" views", ","}, "")
	if santized != "15309" {
		t.Error("Expected 15309 but got", santized)
	}
}
