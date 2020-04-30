# grawl

Simple web crawler for learning.

[![Build Status](https://travis-ci.com/greytabby/grawl.svg?branch=master)](https://travis-ci.com/greytabby/grawl)
![Go](https://github.com/greytabby/grawl/workflows/Go/badge.svg)

## Example

```go
package main

import (
	"fmt"

	"github.com/greytabby/grawl/crawler"
)

func main() {
	URL := "https://hub.docker.com"

	// Define hosts which crawler can visit.
	lr := crawler.NewLimitRule()
	lr.AddAllowedHosts("hub.docker.com")

	c := crawler.NewCrawlerWithLimitRule(URL, 2, lr)
	// Register resonse and error handling functions.
	c.OnVisited(func(cr *crawler.CrawlResult) {
		fmt.Println("Visited:", cr.URL)
	})
	c.OnError(func(err error) {
		fmt.Println("Error:", err.Error())
	})

	// UseHeadlessChrome on request.
	c.UseHeadlessChrome()

	c.SetParallelism(3)

	// Start crawling.
	c.Crawl()
}
```

Output

```text
Visited: https://hub.docker.com
Error: Forbidden host: blog.docker.com
Error: Forbidden host: docs.docker.com
Error: Forbidden host: www.docker.com
Error: Already visited: https://hub.docker.com/
Error: Already visited: https://hub.docker.com/search
Visited: https://hub.docker.com/signup
Error: Invalid URL: 
Error: Forbidden host: docs.docker.com
Error: Forbidden host: www.docker.com
Visited: https://hub.docker.com/_/busybox
Visited: https://hub.docker.com/search
Visited: https://hub.docker.com/_/node
Error: Forbidden host: twitter.com
Error: Forbidden host: www.docker.com
Error: Forbidden host: www.facebook.com
Error: Forbidden host: www.docker.com
Error: Forbidden host: www.docker.com
Error: Forbidden host: www.youtube.com
Error: Forbidden host: www.docker.com
Error: Already visited: https://hub.docker.com/signup
.
.
.

```
