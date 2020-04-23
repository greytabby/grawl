package crawler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	f, err := os.Open("testdata/crawl.html")
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

func TestCrawl(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()
	baseURL := ts.URL
	depth := 2
	buf := bytes.NewBuffer([]byte{})
	crawler := NewCrawler(baseURL, depth, buf)
	crawler.Crawl()

	want := strings.Join([]string{
		baseURL,
		baseURL + "/image/test1.png",
		baseURL + "/image/test2.jpg",
		"",
	}, "\n")
	got := buf.String()

	assert.Equal(t, want, got)
}
