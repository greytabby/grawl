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

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/greytabby/grawl/scrape"
)

type Crawler struct {
	baseRawURL string
	baseURL    *url.URL
	MaxDepth   int
	w          io.Writer
	limitRule  *LimitRule
}

var defaultLimitRule = NewLimitRule()

// NewCrawler returns `*Crawler`.
func NewCrawler(URL string, maxDepth int, w io.Writer) *Crawler {
	u, _ := url.Parse(URL)
	return &Crawler{baseRawURL: URL, baseURL: u, MaxDepth: maxDepth, w: w, limitRule: defaultLimitRule}
}

// NewCrawlerWithLimitRule returns `*Crawler` with LimitRule.
func NewCrawlerWithLimitRule(URL string, maxDepth int, w io.Writer, limitRule *LimitRule) *Crawler {
	u, _ := url.Parse(URL)
	return &Crawler{baseRawURL: URL, baseURL: u, MaxDepth: maxDepth, w: w, limitRule: limitRule}
}

// Crawl start crawling
func (c *Crawler) Crawl() {
	c.crawl(c.baseRawURL, 1)
}

func (c *Crawler) crawl(URL string, depth int) {
	if !c.canVisit(URL, depth) {
		return
	}

	fmt.Fprintf(c.w, "%s\n", URL)
	body, err := c.fetch(URL)
	if err != nil {
		return
	}

	links, err := c.extractLinks(body)
	if err != nil {
		return
	}

	u, _ := url.Parse(URL)
	for _, l := range links {
		u := c.fixURL(u, l)
		c.crawl(u, depth+1)
	}
}

func (c *Crawler) canVisit(URL string, depth int) bool {
	if depth > c.MaxDepth {
		return false
	}

	if !isValidURL(URL) {
		return false
	}

	u, err := url.Parse(URL)
	if err != nil {
		return false
	}

	if !c.limitRule.IsAllow(u) {
		return false
	}
	return true
}

func isValidURL(URL string) bool {
	u, err := url.Parse(URL)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}

func (c *Crawler) fetch(URL string) (body []byte, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Crawler) extractLinks(body []byte) (links []string, err error) {
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

func (c *Crawler) fixURL(currentURL *url.URL, nextURL string) string {
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
