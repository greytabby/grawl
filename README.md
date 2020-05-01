# grawl

Simple web crawler for learning.

[![Build Status](https://travis-ci.com/greytabby/grawl.svg?branch=master)](https://travis-ci.com/greytabby/grawl)
![Go](https://github.com/greytabby/grawl/workflows/Go/badge.svg)

## Usage

```Text
Usage of Grawl:
  -allowed_hosts string
        Accessibel hosts. Use comma to specify multiple hosts
  -depth int
        Limit number of follow links on crawling (default 1)
  -headless_chrome
        Use headless chrome on crawling
  -output_dir string
        Directory name for saving crawl result
  -parallelism int
        Number of parallel execution of crawler (default 5)
  -site string
        Site to crawl
  -v    show version
```

## Example

### Command

```sh
./Grawl -site "https://hub.docker.com" \
-allowed_hosts hub.docker.com \
-depth 2 \
-headless_chrome \
-parallelism 10 \
-output_dir /tmp/dockerhub
```

### Output

console log

```text
Grawl 2020/05/01 15:44:53 Output base directory: /tmp/dockerhub
Grawl 2020/05/01 15:44:53 Crawling site: https://hub.docker.com
Grawl 2020/05/01 15:44:53 Crawling max depth: 2
Grawl 2020/05/01 15:44:53 Start Crawling...
Grawl 2020/05/01 15:44:57 Visited: https://hub.docker.com
Grawl 2020/05/01 15:44:57 Forbidden host: blog.docker.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.docker.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.docker.com
Grawl 2020/05/01 15:44:57 Already visited: https://hub.docker.com/
Grawl 2020/05/01 15:44:57 Forbidden host: www.docker.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.docker.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.linkedin.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.youtube.com
Grawl 2020/05/01 15:44:57 Forbidden host: www.docker.com
Grawl 2020/05/01 15:44:57 Already visited: https://hub.docker.com/search
Grawl 2020/05/01 15:45:06 Visited: https://hub.docker.com/signup
Grawl 2020/05/01 15:45:06 Visited: https://hub.docker.com/_/redis
Grawl 2020/05/01 15:45:07 Visited: https://hub.docker.com/_/ubuntu
Grawl 2020/05/01 15:45:07 Visited: https://hub.docker.com/_/nginx
Grawl 2020/05/01 15:45:07 Visited: https://hub.docker.com/_/postgres
Grawl 2020/05/01 15:45:07 Visited: https://hub.docker.com/_/node
Grawl 2020/05/01 15:45:07 Visited: https://hub.docker.com/_/alpine
Grawl 2020/05/01 15:45:08 Visited: https://hub.docker.com/_/couchbase
Grawl 2020/05/01 15:45:08 Already visited: https://hub.docker.com/signup
Grawl 2020/05/01 15:45:08 Already visited: https://hub.docker.com/search
Grawl 2020/05/01 15:45:08 Invalid URL: 
Grawl 2020/05/01 15:45:08 Invalid URL: 
.
.
.

```

output directory tree

```text
/tmp/dockerhub
└── hub_docker_com
    ├── _
    │   ├── aerospike
    │   │   └── index.html
    │   ├── alpine
    │   │   └── index.html
    │   ├── busybox
    │   │   └── index.html
    │   ├── couchbase
    │   │   └── index.html
    │   ├── golang
    │   │   └── index.html
    │   ├── hello-world
    │   │   └── index.html
    │   ├── mongo
    │   │   └── index.html
    │   ├── mysql
    │   │   └── index.html
    │   ├── nginx
    │   │   └── index.html
    │   ├── node
    │   │   └── index.html
    │   ├── postgres
    │   │   └── index.html
    │   ├── redis
    │   │   └── index.html
    │   ├── registry
    │   │   └── index.html
    │   ├── traefik
    │   │   └── index.html
    │   └── ubuntu
    │       └── index.html
    ├── index.html
    ├── search
    │   └── index.html
    └── signup
        └── index.html
```

## docker-compose

```yml
version: '2'
services: 
  grawl:
    image: greytabby/grawl:latest
    volumes:
      - .:/result/dockerhub
    environment: 
      # SITE is url which grawl first visit
      - SITE=https://hub.docker.com/
      # ALLOWED_HOSTS are acceessible hosts. Use comma to specify multiple hosts
      - ALLOWED_HOSTS=hub.docker.com
      # DEPTH is limit nuber of follow links on crawling
      - DEPTH=2
      # PARALLELISM is number of parallel execution of crawler
      - PARALLELISM=5
      # HEADLESS_CHROME use headless chreme on crawling
      # If you do not use headless chrome, delete this variable
      - HEADLESS_CHROME=
      # OUTPUT_DIR is directory name for saving crawl result
      - OUTPUT_DIR=/result/dockerhub
```
