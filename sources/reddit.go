package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RedditComment is a comment returned from Pushshift.
type RedditComment struct {
	Body string
}

// RedditCommentData is data returned from a Pushshift comment query.
type RedditCommentData struct {
	Data []RedditComment
}

// RedditCommentFetcher is a thing that can fetch comments from REddit.
type RedditCommentFetcher interface {
	// GetComments retrieves comments.
	GetComments() ([]byte, error)
}

// RedditClient is a simple Reddit client.
type RedditClient struct {
	// PushshiftBaseURL is the base URL for the Pushshift query
	PushshiftBaseURL string

	// PushshiftQueryParameters are query params to provide to
	// Pushshift.
	PushshiftQueryParameters string
}

// GetComments gets comments from Pushshift.
func (c *RedditClient) GetComments() ([]byte, error) {
	url := fmt.Sprintf("%s?%s", c.PushshiftBaseURL, c.PushshiftQueryParameters)
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

func NewRedditClient(qp string) *RedditClient {
	defaultParams := strings.Join([]string{"q=\"devops agile\"",
		"subreddit=devops",
		"fields=\"author,body\"",
		"metadata=false",
	}, "&")
	var params string
	if qp != "" {
		params = qp
	} else {
		params = defaultParams
	}
	return &RedditClient{
		PushshiftBaseURL:         "https://api.pushshift.io/reddit/comment/search",
		PushshiftQueryParameters: params,
	}
}

func commentListFromJSON(comments []byte) ([]string, error) {
	var r RedditCommentData
	err := json.Unmarshal(comments, &r)
	if err != nil {
		return []string{}, err
	}
	fmt.Printf("comments: %d\n", len(r.Data))
	var cl []string
	for _, d := range r.Data {
		cl = append(cl, d.Body)
	}
	return cl, nil
}

func GenerateWordCountReddit(client RedditCommentFetcher) (map[string]int, error) {
	counts := map[string]int{}
	comments, err := client.GetComments()
	if err != nil {
		return counts, err
	}
	commentList, err := commentListFromJSON(comments)
	if err != nil {
		return counts, err
	}
	return countWordsInComments(commentList), nil
}
