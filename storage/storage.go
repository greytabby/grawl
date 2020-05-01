package storage

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/greytabby/grawl/crawler"
)

type FileStorage struct {
	BaseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{baseDir}
}

func (fs *FileStorage) Save(cr *crawler.CrawlResult) error {
	path := fs.urlToFilepath(cr.URL)
	err := os.MkdirAll(filepath.Dir(filepath.Clean(path)), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(cr.Body), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) urlToFilepath(URL *url.URL) string {
	if URL.Host == "" {
		return ""
	}

	host := strings.ReplaceAll(URL.Host, ":", "-") // port number
	host = strings.ReplaceAll(host, ".", "_")      // host
	path := filepath.Join(fs.BaseDir, host, URL.Path)
	if filepath.Ext(filepath.Base(path)) == "" {
		path = filepath.Join(path, "index.html")
	}
	return path
}
