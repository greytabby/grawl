package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/greytabby/grawl/fetcher"
	"github.com/greytabby/grawl/scrape"
)

type (
	Crawler struct {
		baseRawURL       string
		maxDepth         int
		fetcher          Fetcher
		limitRule        *LimitRule
		parallelism      chan struct{}
		visitCallbacks   []VisitCallback
		visitedCallbacks []VisitedCallback
		errorCallbacks   []ErrorCallback
		set              map[string]bool
		mux              sync.RWMutex
		wg               sync.WaitGroup
	}

	CrawlResult struct {
		URL   *url.URL
		Body  string
		Links []string
	}
)

var (
	defaultLimitRule   = NewLimitRule()
	defaultParallelism = 5
)

var (
	// ErrInvalidURL is the error thrown if visiting URL
	// is invalid format.
	ErrInvalidURL = errors.New("Invalid URL")
	// ErrForbidden is the error thrown if the url is not allowed to visit
	ErrForbidden = errors.New("Forbidden")
	// ErrAlreadyVisitedDomain is the error for already visited URL
	ErrAlreadyVisited = errors.New("Already visited")
)

type (
	// VisitCallback is a type of alias for OnVisit callback functions
	VisitCallback func(response []byte)
	// VisitedCallback is a type of alias for OnVisited callback functions
	VisitedCallback func(*CrawlResult)
	// ErrorCallback is a type of alias for OnError callback functions
	ErrorCallback func(error)
)

// Fetcher sends GET request to the given URL and
// returns response body
type Fetcher interface {
	Fetch(URL string) (body []byte, err error)
}

// NewCrawler returns `*Crawler`.
func NewCrawler(URL string, maxDepth int) *Crawler {
	return &Crawler{
		baseRawURL:       URL,
		maxDepth:         maxDepth,
		fetcher:          new(fetcher.DefaultFetcher),
		limitRule:        defaultLimitRule,
		parallelism:      make(chan struct{}, defaultParallelism),
		visitedCallbacks: []VisitedCallback{},
		set:              map[string]bool{},
	}
}

// NewCrawlerWithLimitRule returns `*Crawler` with LimitRule.
func NewCrawlerWithLimitRule(URL string, maxDepth int, limitRule *LimitRule) *Crawler {
	c := NewCrawler(URL, maxDepth)
	c.limitRule = limitRule
	return c
}

// UseHeadlessChrome use headless chrome at the time of request.
// By default, using `http.Get()`.
func (c *Crawler) UseHeadlessChrome() {
	c.fetcher = new(fetcher.HeadlessChrome)
}

// SetParallelism set limit of crawling parallelism.
// By default, parallelism is 5.
func (c *Crawler) SetParallelism(n int) {
	c.parallelism = make(chan struct{}, n)
}

// OnVisit register a function. Function will be executed on visiting web site.
func (c *Crawler) OnVisit(f VisitCallback) {
	c.visitCallbacks = append(c.visitCallbacks, f)
}

// OnVisited register a function. Function will be executed on after
// visit web site.
func (c *Crawler) OnVisited(f VisitedCallback) {
	c.visitedCallbacks = append(c.visitedCallbacks, f)
}

// OnError register a function. Function will be executed on error occured
func (c *Crawler) OnError(f ErrorCallback) {
	c.errorCallbacks = append(c.errorCallbacks, f)
}

// Crawl start crawling
func (c *Crawler) Crawl() {
	c.wg.Add(1)
	go c.crawl(c.baseRawURL, 1)
	c.wg.Wait()
}

func (c *Crawler) crawl(rawURL string, depth int) {
	defer c.wg.Done()
	c.parallelism <- struct{}{}
	defer func() {
		<-c.parallelism
	}()
	if depth > c.maxDepth {
		return
	}

	URL, err := url.Parse(rawURL)
	if err != nil {
		c.handleErrorCallback(err)
		return
	}
	if err = c.canVisit(URL); err != nil {
		c.handleErrorCallback(err)
		return
	}

	cr, err := c.visit(URL)
	if err != nil {
		c.handleErrorCallback(err)
		return
	}

	for _, link := range cr.Links {
		nextRawURL := fixURL(URL, link)
		c.wg.Add(1)
		go c.crawl(nextRawURL, depth+1)
	}
}

func (c *Crawler) canVisit(URL *url.URL) error {
	if !isValidURL(URL) {
		return fmt.Errorf("%v: %s", ErrInvalidURL, URL)
	}

	if !c.limitRule.IsAllow(URL) {
		return fmt.Errorf("%v: %s", ErrForbidden, URL.String())
	}

	if c.hasVisited(toNoneQueryAndFragmentURL(URL)) {
		return fmt.Errorf("%v: %s", ErrAlreadyVisited, URL)
	}

	return nil
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

func (c *Crawler) visit(URL *url.URL) (*CrawlResult, error) {
	u := toNoneQueryAndFragmentURL(URL)
	c.setVisit(u)
	body, err := c.fetcher.Fetch(URL.String())
	if err != nil {
		return nil, err
	}
	c.handleVisitCallback(body)

	links, err := extractLinks(body)
	if err != nil {
		return nil, err
	}
	cr := &CrawlResult{URL, string(body), links}
	c.handleVisitedCallback(cr)
	return cr, nil
}

func toNoneQueryAndFragmentURL(URL *url.URL) string {
	u := URL.Scheme + "://" + URL.Host + URL.Path
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

func (c *Crawler) handleVisitCallback(response []byte) {
	for _, f := range c.visitCallbacks {
		f(response)
	}
}

func (c *Crawler) handleVisitedCallback(cr *CrawlResult) {
	for _, f := range c.visitedCallbacks {
		f(cr)
	}
}

func (c *Crawler) handleErrorCallback(err error) {
	for _, f := range c.errorCallbacks {
		f(err)
	}
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
