package crawler

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedHost(t *testing.T) {
	testCases := []struct {
		host string
		want bool
	}{
		{"test.com", true},
		{"www.example.com", true},
		{"neco.net", true},
		{"www.test.com", false},
		{"example.com", false},
		{"", false},
	}

	limitRule := NewLimitRule()
	limitRule.AddAllowedHosts("test.com", "www.example.com", "neco.net")
	for _, tt := range testCases {
		got := limitRule.isAllowedHost(tt.host)
		assert.Equal(t, tt.want, got)
	}
}

func TestIsAllowedHostWithNoneAllowedHosts(t *testing.T) {
	testCases := []struct {
		host string
		want bool
	}{
		{"test.com", true},
		{"www.example.com", true},
		{"neco.net", true},
		{"www.test.com", true},
		{"example.com", true},
		{"", true},
	}

	limitRule := NewLimitRule()
	for _, tt := range testCases {
		got := limitRule.isAllowedHost(tt.host)
		assert.Equal(t, tt.want, got)
	}
}

func TestIsAllow(t *testing.T) {
	testCases := []struct {
		host string
		want bool
	}{
		{"test.com", true},
		{"www.example.com", true},
		{"neco.net", true},
		{"www.test.com", false},
		{"example.com", false},
		{"", false},
	}

	limitRule := NewLimitRule()
	limitRule.AddAllowedHosts("test.com", "www.example.com", "neco.net")
	baseURL := "http://"
	for _, tt := range testCases {
		requestURL, err := url.Parse(baseURL + tt.host)
		assert.NoError(t, err)
		got := limitRule.IsAllow(requestURL)
		assert.Equal(t, tt.want, got)
	}
}
