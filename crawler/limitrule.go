package crawler

import (
	"net/url"
	"regexp"
)

type LimitRule struct {
	// AllowedHosts define accessible hosts.
	// When AllowedHosts is empty, all hosts are allowed.
	AllowedHosts []string
	AllowedUrls  []*regexp.Regexp
}

// NewLimitRule returns empty LimitRule.
// You can add rule to it to call AddAllowedHosts
func NewLimitRule() *LimitRule {
	lr := new(LimitRule)
	lr.AllowedHosts = []string{}
	return lr
}

// IsAllow returns true if requestURL is no limit to crawl.
func (lr *LimitRule) IsAllow(requestURL *url.URL) bool {
	if lr.isAllowedHost(requestURL.Host) {
		if len(lr.AllowedUrls) == 0 {
			return true
		}

		for _, allowedURL := range lr.AllowedUrls {
			if allowedURL.Match([]byte(requestURL.String())) {
				return true
			}
		}
	}

	return false
}

// AddAllowedHosts add rule define accessible hosts.
func (lr *LimitRule) AddAllowedHosts(hosts ...string) {
	lr.AllowedHosts = append(lr.AllowedHosts, hosts...)
}

func (lr *LimitRule) AddAllowedUrls(re *regexp.Regexp) {
	lr.AllowedUrls = append(lr.AllowedUrls, re)
}

func (lr *LimitRule) isAllowedHost(host string) bool {
	if len(lr.AllowedHosts) == 0 {
		return true
	}

	for _, h := range lr.AllowedHosts {
		if h == host {
			return true
		}
	}
	return false
}
