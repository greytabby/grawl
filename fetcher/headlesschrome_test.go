package fetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeHTML(rawHTML string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, rawHTML)
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(writeHTML(`<!doctype html>
<html>
<head>
  <link rel="stylesheet" href="test.css">
</head>
<body>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>
</body>
</html>`))
	return ts
}

func TestHeadlessChromeFetch(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	hc := new(HeadlessChrome)
	body, err := hc.Fetch(ts.URL)
	assert.NoError(t, err)

	want := []byte(`<head>
  <link rel="stylesheet" href="test.css">
</head>
<body>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>
  <div id="content">the content</div>

</body>`)

	assert.Equal(t, want, body)
}