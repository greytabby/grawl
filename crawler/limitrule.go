package crawler

import "net/url"

type LimitRule struct {
	// AllowedHosts define accessible hosts.
	// When AllowedHosts is empty, all hosts are allowed.
	AllowedHosts []string
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
	if !lr.isAllowedHost(requestURL.Host) {
		return false
	}
	return true
}

// AddAllowedHosts add rule define accessible hosts.
func (lr *LimitRule) AddAllowedHosts(hosts ...string) {
	lr.AllowedHosts = append(lr.AllowedHosts, hosts...)
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
