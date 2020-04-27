package fetcher

import (
	"io/ioutil"
	"net/http"
)

type DefaultFetcher struct{}

func (df *DefaultFetcher) Fetch(URL string) (body []byte, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
