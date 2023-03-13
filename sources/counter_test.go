package sources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterCounts(t *testing.T) {
	sentences := `the quick
the quick fox
The quick fox jumped
the quick fox jumped over
The quick fox jumped over the
The quick fox jumped over the lazy
the quick fox jumped over the lazy dog`
	want := map[string]int{
		"lazy":   2,
		"jumped": 5,
		"fox":    6,
		"quick":  7,
	}
	got := countWordsInComments(strings.Split(sentences, "\n"))
	assert.Equal(t, want, got)
}
