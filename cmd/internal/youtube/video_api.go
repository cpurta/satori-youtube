package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	body, _ := ioutil.ReadAll(res.Body)

	var videoResponse VideoListResponse
	err = json.Unmarshal(body, &videoResponse)

	return &videoResponse, err
}
