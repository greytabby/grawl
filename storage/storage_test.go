package storage

import (
	"io/ioutil"
	"testing"
	"net/url"
	"os"

	"github.com/greytabby/grawl/crawler"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	visited := []string{
		"https://test.com/",
		"https://test.com/file1",
		"https://test.com/users/file2",
		"https://test.com:8080/file3",
		"https://docs.test.com/file4",
		"https://diffrenthost.aaa.bbb/",
		"https://diffrenthost.aaa.bbb/file5",
	}

	urls := make([]*url.URL, len(visited))
	crawlResults := make([]*crawler.CrawlResult, len(visited))
	for i, v := range visited {
		urls[i], _ = url.Parse(v)
		crawlResults[i] = &crawler.CrawlResult{URL: urls[i], Body: v, Links: []string{}}
	}

	tempDir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(tempDir) // tier down
	storage := NewFileStorage(tempDir)
	for _, v := range crawlResults {
		err = storage.Save(v)
		assert.NoError(t, err)
		path := storage.urlToFilepath(v.URL)
		assert.FileExists(t, path)
		t.Log("Path:", path)
	}
}