# grawl

Simple web crawler for learning.

[![Build Status](https://travis-ci.com/greytabby/grawl.svg?branch=master)](https://travis-ci.com/greytabby/grawl)
![Go](https://github.com/greytabby/grawl/workflows/Go/badge.svg)

## Example

```go
package main

import (
	"os"

	"github.com/greytabby/grawl/crawler"
)

func main() {
	URL := "https://hub.docker.com/"
	depth := 2
	output := os.Stdout
	lr := crawler.NewLimitRule()
	lr.AddAllowedHosts("hub.docker.com")
	c := crawler.NewCrawlerWithLimitRule(URL, depth, output, lr)
	c.UseHeadlessChrome()
	c.SetParallelism(10)
	c.Crawl()
}
```

Output

```text
https://hub.docker.com/
https://hub.docker.com/search
https://hub.docker.com/signup
https://hub.docker.com/_/nginx
https://hub.docker.com/_/mongo
https://hub.docker.com/_/alpine
https://hub.docker.com/_/node
https://hub.docker.com/_/redis
https://hub.docker.com/_/couchbase
https://hub.docker.com/_/ubuntu
https://hub.docker.com/_/mysql
https://hub.docker.com/_/postgres
https://hub.docker.com/_/traefik
https://hub.docker.com/_/busybox
https://hub.docker.com/_/mariadb
https://hub.docker.com/_/registry
https://hub.docker.com/_/hello-world
https://hub.docker.com/_/docker
.
.
.

```
