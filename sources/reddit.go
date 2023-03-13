package sources

import (
	"encoding/json"
	"fmt"
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
	// PushshiftQueryString is the query string to use to fetch comments.
	PushshiftQueryString string
}

// GetComments gets comments from Pushshift.
func (c *RedditClient) GetComments() ([]byte, error) {
	return []byte("wip"), nil
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
