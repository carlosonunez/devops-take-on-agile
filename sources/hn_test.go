package sources

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockHNClient struct {
	m sync.Mutex
}

func (m *MockHNClient) GetPageCount() (int, error) {
	return 5, nil
}

func (m *MockHNClient) GetCommentsByPage(page int) ([]string, error) {
	comments, err := os.ReadFile(fmt.Sprintf("./fixtures/hn_comments_page_%d.json", page))
	if err != nil {
		return []string{}, fmt.Errorf("couldn't open comments: %s", err.Error())
	}
	return getCommentsList(comments)
}

func TestGettingCommentsFromHNPayload(t *testing.T) {
	body, err := os.ReadFile(fmt.Sprintf("./fixtures/hn_comments_page_1.json"))
	if err != nil {
		t.Logf("Couldn't read comments: %s", err.Error())
		t.Fail()
	}
	want := []string{
		"SEEKING WORK | Norway | Remote | Contract (full-time or part-time)<p>I am a software engineer with over eight years of experience and a background in information security. I prefer backend work, and I have experience in a variety of projects and technologies, including data integration, Intel SGX, consensus protocols, REST APIs and web development. I am familiar with agile methodologies, and I have worked as scrum master and in DevOps roles.<p>I learn things quickly, greatly enjoy solving challenging problems and pride myself in delivering custom solutions that address your needs in the best possible way, from their design and architecture to their day-to-day performance.<p>In addition to doing development, I believe I am the most helpful with system design and architecture, including APIs; coaching and mentoring and code reviews.<p>Technologies: C, C++, CSS, Docker, ES6+, express.js, Intel SGX, Java, JavaScript, Kotlin, LDAP, Linux, Neo4j, nginx, Node.js, PHP, PL&#x2F;SQL, Postfix, React, TypeScript, Xen, (X)HTML5 and more. Since Recently a W3C member, so if you&#x27;ve work with &#x2F; around web standards, I&#x27;d be stoked to help!<p>I am available for full-time and part-time engagements as well as contract work.<p>Location: Trøndelag, Norway<p>Remote: Yes (remote only, unless within Trøndelag or occasional meetups within Scandinavia)<p>Willing to relocate: No<p>Résumé&#x2F;CV: <a href=\"https:&#x2F;&#x2F;xcrty.link&#x2F;vn1s\" rel=\"nofollow\">https:&#x2F;&#x2F;xcrty.link&#x2F;vn1s</a><p>Email: hn-u5cgNWJM [at] protonmail.com<p>GitHub: <a href=\"https:&#x2F;&#x2F;github.com&#x2F;corrideat\">https:&#x2F;&#x2F;github.com&#x2F;corrideat</a>",
	}
	got, err := getCommentsList(body)
	fmt.Printf("===> num [test]: %d\n", len(got))
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// The Hacker News source uses Algolia since HN doesn't have a "real" API.
func TestHNWordCountGeneratorPass(t *testing.T) {
	if os.Getenv("INTEGRATION") == "1" {
		t.Logf("skipping unit test")
		t.Skip()
	}
	counts, err := os.ReadFile("./fixtures/hn_expected_word_count.txt")
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
	got, err := GenerateWordCountHackerNews(&MockHNClient{})
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestHNWordCountGeneratorIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Logf("skipping integration test")
		t.Skip()
	}
	hc := NewHNClient(strings.Join([]string{
		"query=devops%20agile",
		"tags=(story_35096366,story_34983765,story_34983766,story_34983765,story_34886374)",
	}, "&"))
	counts, err := os.ReadFile("./fixtures/hn_expected_word_count.txt")
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
	got, err := GenerateWordCountHackerNews(hc)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
