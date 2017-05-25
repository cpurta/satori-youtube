package youtube

import (
	"fmt"
	"testing"
)

func TestScrapeSnippetData(t *testing.T) {
	snippet, err := ScrapeSnippetData("https://www.youtube.com/watch?v=i8YRtIHRIv0")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Title:", snippet.Title)
	fmt.Println("Description:", snippet.Description)
	fmt.Println("Channel Title:", snippet.ChannelTitle)
	fmt.Println("Published at:", snippet.PublishedAt)
}

func TestScrapeStatisticsData(t *testing.T) {
	stats, err := ScrapeStatisticsData("https://www.youtube.com/watch?v=i8YRtIHRIv0")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("View count:", stats.ViewCount)
	fmt.Println("Like count", stats.LikeCount)
	fmt.Println("Dislike count", stats.DislikeCount)
	fmt.Println("Comment count", stats.CommentCount)
}
