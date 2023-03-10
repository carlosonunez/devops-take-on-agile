package sources

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// The reddit source uses the Pushshift API to retrieve comments instead of
// Reddit's API directly. We're using Pushshift because Reddit doesn't make
// comments publicly accessible.
func TestRedditWordCountGenerator(t *testing.T) {
	want, err := os.ReadFile("./fixtures/reddit_expected_word_count.txt")
	if err != nil {
		t.Logf("couldn't open expected data: %s", err.Error())
		t.FailNow()
	}
	got := GenerateWordCountReddit()
	assert.Equal(want, got)
}
