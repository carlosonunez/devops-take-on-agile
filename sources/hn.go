package sources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"
)

const HNConcurrencyLimit = 5

// HNComment is a comment returned from Algolia.
type HNComment struct {
	CommentText string `json:"comment_text"`
}

// HNCommentData is data returned from a Algolia comment query.
type HNCommentData struct {
	Hits     []HNComment
	NumPages int `json:"numPages"`
}

// HNCommentFetcher is a thing that can fetch comments from REddit.
type HNCommentFetcher interface {
	// GetCommentsByPage retrieves comments.
	GetCommentsByPage(int) ([]string, error)

	// GetPageCount returns the number of pages in the set of results.
	GetPageCount() (int, error)
}

// HNClient is a simple HN client.
type HNClient struct {
	// AlgoliaBaseURL is the base URL for the Algolia query
	AlgoliaBaseURL string

	// AlgoliaQueryParameters are query params to provide to
	// Algolia.
	AlgoliaQueryParameters string
}

// GetPageCount returns the number of pages in a query result set.
func (c *HNClient) GetPageCount() (int, error) {
	url := fmt.Sprintf("%s?%s&hitsPerPage=20", c.AlgoliaBaseURL, c.AlgoliaQueryParameters)
	body, err := getFromAlgolia(url)
	if err != nil {
		return 0, err
	}
	hnData := HNCommentData{}
	if err := json.Unmarshal(body, &hnData); err != nil {
		return 0, err
	}
	return hnData.NumPages, nil
}

// GetCommentsByPage gets comments from Algolia.
func (c *HNClient) GetCommentsByPage(page int) ([]string, error) {
	url := fmt.Sprintf("%s?%s&page=%d", c.AlgoliaBaseURL, c.AlgoliaQueryParameters, page)
	body, err := getFromAlgolia(url)
	if err != nil {
		return []string{}, err
	}
	return getCommentsList(body)
}

func NewHNClient(qp string) *HNClient {
	defaultParams := strings.Join([]string{"query=\"devops agile\"",
		"tags=comment",
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
	commentList, err := getCommentsParallel(client)
	if err != nil {
		return counts, err
	}
	if err != nil {
		return counts, err
	}
	return countWordsInComments(commentList), nil
}

func getCommentsParallel(c HNCommentFetcher) ([]string, error) {
	numPages, err := c.GetPageCount()
	comments := []string{}
	if err != nil {
		return comments, err
	}
	var g errgroup.Group
	g.SetLimit(HNConcurrencyLimit)
	var results = make([][]string, numPages)
	for page := 0; page < numPages; page++ {
		p := page
		g.Go(func() error {
			fmt.Printf("Page: %d\n", p)
			commentList, err := c.GetCommentsByPage(p)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			results[p] = commentList
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return comments, err
	}
	for _, r := range results {
		for _, c := range r {
			comments = append(comments, c)
		}
	}
	return comments, nil
}

func getFromAlgolia(url string) ([]byte, error) {
	var body []byte
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
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func getCommentsList(body []byte) ([]string, error) {
	hnData := HNCommentData{}
	if err := json.Unmarshal(body, &hnData); err != nil {
		return []string{}, err
	}
	var comments = make([]string, len(hnData.Hits), len(hnData.Hits))
	for i := 0; i < len(hnData.Hits); i++ {
		comments[i] = hnData.Hits[i].CommentText
	}
	return comments, nil
}
