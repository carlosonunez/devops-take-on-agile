package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HNComment is a comment returned from Algolia.
type HNComment struct {
	CommentText string `json:"comment_text"`
}

// HNCommentData is data returned from a Algolia comment query.
type HNCommentData struct {
	Hits []HNComment
}

// HNCommentFetcher is a thing that can fetch comments from REddit.
type HNCommentFetcher interface {
	// GetComments retrieves comments.
	GetComments() ([]byte, error)
}

// HNClient is a simple HN client.
type HNClient struct {
	// AlgoliaBaseURL is the base URL for the Algolia query
	AlgoliaBaseURL string

	// AlgoliaQueryParameters are query params to provide to
	// Algolia.
	AlgoliaQueryParameters string
}

// GetComments gets comments from Algolia.
func (c *HNClient) GetComments() ([]byte, error) {
	url := fmt.Sprintf("%s?%s", c.AlgoliaBaseURL, c.AlgoliaQueryParameters)
	fmt.Printf("Getting comments from %s\n", url)
	resp, err := http.Get(url)
	resp.Header.Set("User-Agent", "curl/7.85.0")
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("expected 200; got %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	fmt.Printf("We got: %s", body)
	return body, nil
}

func NewHNClient(qp string) *HNClient {
	defaultParams := strings.Join([]string{"query=\"devops agile\"",
		"tags=comment",
		"hitsPerPage=1000000",
		"numericFilters=created_at_i>1609480800",
	}, "&")
	var params string
	if qp != "" {
		params = qp
	} else {
		params = defaultParams
	}
	return &HNClient{
		AlgoliaBaseURL:         "https://hn.algolia.com/api/v1/search_by_date",
		AlgoliaQueryParameters: params,
	}
}

func GenerateWordCountHackerNews(client HNCommentFetcher) (map[string]int, error) {
	counts := map[string]int{}
	comments, err := client.GetComments()
	if err != nil {
		return counts, err
	}
	commentList, err := commentListFromJSONHN(comments)
	if err != nil {
		return counts, err
	}
	return countWordsInComments(commentList), nil
}

func commentListFromJSONHN(comments []byte) ([]string, error) {
	var r HNCommentData
	err := json.Unmarshal(comments, &r)
	if err != nil {
		return []string{}, err
	}
	var cl []string
	for _, d := range r.Hits {
		cl = append(cl, d.CommentText)
	}
	return cl, nil
}
