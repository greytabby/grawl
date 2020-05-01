package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, testfile string) *httptest.Server {
	t.Helper()

	f, err := os.Open(testfile)
	assert.NoError(t, err)
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			fmt.Fprint(w, "test")
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, string(content))
	})

	ts := httptest.NewServer(handler)
	return ts
}

func confirmCrawlResult(t *testing.T, want, got []string) bool {
	wantl := len(want)
	count := 0
	for _, g := range got {
		for _, w := range want {
			if g == w {
				count++
				break
			}
		}
	}
	return wantl == count
}

func TestCrawl(t *testing.T) {
	ts := newTestServer(t, "testdata/crawl.html")
	defer ts.Close()
	baseURL := ts.URL
	depth := 2
	c := NewCrawler(baseURL, depth)
	got := make([]string, 0)
	c.OnVisited(func(cr *CrawlResult) {
		got = append(got, cr.URL.String())
	})
	c.Crawl()

	want := []string{
		baseURL,
		baseURL + "/image/test1.png",
		baseURL + "/image/test2.jpg",
		baseURL + "/relative/test3",
	}

	if !confirmCrawlResult(t, want, got) {
		t.Error("want:", want, "got:", got)
	}
}

func TestCrawlWithHostsLimit(t *testing.T) {
	ts := newTestServer(t, "testdata/crawl-with-host-limit.html")
	defer ts.Close()
	URL, _ := url.Parse(ts.URL)
	limitRule := NewLimitRule()
	limitRule.AddAllowedHosts(URL.Host)
	c := NewCrawlerWithLimitRule(ts.URL, 2, limitRule)
	got := make([]string, 0)
	c.OnVisited(func(cr *CrawlResult) {
		got = append(got, cr.URL.String())
	})
	c.Crawl()

	want := []string{
		ts.URL,
		ts.URL + "/image/test1.png",
		ts.URL + "/image/test2.jpg",
	}
	if !confirmCrawlResult(t, want, got) {
		t.Error("want:", want, "got:", got)
	}
}

func TestCrawlDontVisitSameURL(t *testing.T) {
	ts := newTestServer(t, "testdata/crawl-dont-visit-same-url.html")
	defer ts.Close()
	got := make([]string, 0)
	c := NewCrawler(ts.URL, 2)
	c.OnVisited(func(cr *CrawlResult) {
		got = append(got, cr.URL.String())
	})
	c.Crawl()

	want := []string{
		ts.URL,
		ts.URL + "/image/test1",
		ts.URL + "/image/test2",
	}
	if !confirmCrawlResult(t, want, got) {
		t.Error("want:", want, "got:", got)
	}
}

func TestCrawlOnError(t *testing.T) {
	t.Run("Invalid URL", func(t *testing.T) {
		c := NewCrawler("//test.com", 1)
		var err error
		c.OnError(func(e error) {
			err = e
		})
		c.Crawl()
		assert.EqualError(t, err, fmt.Sprintf("%v: %s", ErrInvalidURL, "//test.com"))
	})

	t.Run("Forbidden host", func(t *testing.T) {
		lr := new(LimitRule)
		lr.AddAllowedHosts("test.com")
		c := NewCrawlerWithLimitRule("https://not.allow.com/index.html", 1, lr)
		var err error
		c.OnError(func(e error) {
			err = e
		})
		c.Crawl()
		assert.EqualError(t, err, fmt.Sprintf("%v: %s", ErrForbiddenHost, "not.allow.com"))
	})

	t.Run("Cannot parse url", func(t *testing.T) {
		c := NewCrawler("?>_&ﬁ·)·‚‚ﬁ°ﬁ‚‚ﬁ", 1)
		var err error
		c.OnError(func(e error) {
			err = e
		})
		c.Crawl()
		assert.Error(t, err)
	})

	// TODO: Test html parse error
}
