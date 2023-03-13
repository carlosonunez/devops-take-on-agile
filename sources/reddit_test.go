package sources

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRedditClient struct{}

func (m *MockRedditClient) GetComments() ([]byte, error) {
	comments, err := os.ReadFile("./fixtures/reddit_comments.json")
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't open comments: %s", err.Error())
	}
	return comments, nil
}

// The reddit source uses the Pushshift API to retrieve comments instead of
// Reddit's API directly. We're using Pushshift because Reddit doesn't make
// comments publicly accessible.
func TestRedditWordCountGeneratorPass(t *testing.T) {
	counts, err := os.ReadFile("./fixtures/reddit_expected_word_count.txt")
	if err != nil {
		t.Logf("couldn't open expected data: %s", err.Error())
		t.FailNow()
	}
	want := map[string]int{}
	for _, c := range strings.Split(string(counts), "\n") {
		if len(strings.TrimSpace(c)) == 0 {
			continue
		}
		data := strings.Split(c, " ")
		count, _ := strconv.Atoi(data[0])
		want[data[1]] = count
	}
	got, _ := GenerateWordCountReddit(&MockRedditClient{})
	assert.Equal(t, want, got)
}
