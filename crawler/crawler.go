package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/greytabby/grawl/scrape"
)

type Crawler struct {
	baseRawURL string
	maxDepth   int
	w          io.Writer
	limitRule  *LimitRule
	set        map[string]bool
	mux        sync.RWMutex
}

var defaultLimitRule = NewLimitRule()

// NewCrawler returns `*Crawler`.
func NewCrawler(URL string, maxDepth int, w io.Writer) *Crawler {
	return &Crawler{
		baseRawURL: URL,
		maxDepth:   maxDepth,
		w:          w,
		limitRule:  defaultLimitRule,
		set:        map[string]bool{},
	}
}

// NewCrawlerWithLimitRule returns `*Crawler` with LimitRule.
func NewCrawlerWithLimitRule(URL string, maxDepth int, w io.Writer, limitRule *LimitRule) *Crawler {
	c := NewCrawler(URL, maxDepth, w)
	c.limitRule = limitRule
	return c
}

// Crawl start crawling
func (c *Crawler) Crawl() {
	c.crawl(c.baseRawURL, 1)
}

func (c *Crawler) crawl(rawURL string, depth int) {
	if depth > c.maxDepth {
		return
	}

	URL, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	if !c.canVisit(URL) {
		return
	}

	links, err := c.visit(URL)
	if err != nil {
		return
	}

	for _, link := range links {
		nextRawURL := fixURL(URL, link)
		c.crawl(nextRawURL, depth+1)
	}
}

func (c *Crawler) canVisit(URL *url.URL) bool {
	if !isValidURL(URL) {
		return false
	}

	if !c.limitRule.IsAllow(URL) {
		return false
	}

	if c.hasVisited(toNoneQueryAndFragmentURL(URL)) {
		return false
	}

	return true
}

func isValidURL(URL *url.URL) bool {
	if URL.Scheme != "http" && URL.Scheme != "https" {
		return false
	}

	if URL.Host == "" {
		return false
	}

	return true
}

func (c *Crawler) hasVisited(URL string) bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	visited, _ := c.set[URL]
	return visited
}

func (c *Crawler) setVisit(URL string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.set[URL] = true
}

func (c *Crawler) visit(URL *url.URL) ([]string, error) {
	fmt.Fprintf(c.w, "%s\n", URL.String())
	u := toNoneQueryAndFragmentURL(URL)
	c.setVisit(u)
	body, err := c.fetch(URL.String())
	if err != nil {
		return []string{}, err
	}

	return extractLinks(body)
}

func (c *Crawler) fetch(URL string) (body []byte, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func toNoneQueryAndFragmentURL(URL *url.URL) string {
	u := URL.Scheme + URL.Host + URL.Path
	return strings.TrimRight(u, "/")
}

func extractLinks(body []byte) (links []string, err error) {
	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	anchorNodes := scrape.FindAll(rootNode, scrape.ByTag(atom.A))
	links = make([]string, len(anchorNodes))
	for i, v := range anchorNodes {
		links[i] = scrape.Attr(v, "href")
	}
	return links, nil
}

func fixURL(currentURL *url.URL, nextURL string) string {
	u, err := url.Parse(nextURL)

	if err != nil {
		return ""
	}

	if u.IsAbs() {
		return u.String()
	}

	if u.Host == "" && u.Path != "" {
		if strings.HasPrefix(u.Path, "/") {
			return currentURL.Scheme + "://" + currentURL.Host + path.Join(currentURL.Path, u.Path)
		}
		return currentURL.Scheme + "://" + currentURL.Host + "/" + path.Join(currentURL.Path, u.Path)
	}

	return ""
}
