FROM golang:1.13-alpine as build

WORKDIR /go/src/grawl
COPY . .
RUN apk add make
RUN make ci-build

FROM chromedp/headless-shell:stable
COPY --from=build /bin/grawl /bin/grawl

# SITE is url which grawl first visit
ENV SITE=
# ALLOWED_HOSTS are acceessible hosts. Use comma to specify multiple hosts
ENV ALLOWED_HOSTS=
# DEPTH is limit nuber of follow links on crawling
ENV DEPTH=1
# PARALLELISM is number of parallel execution of crawler
ENV PARALLELISM=5
# HEADLESS_CHROME use headless chreme on crawling
ENV HEADLESS_CHROME=true
# OUTPUT_DIR is directory name for saving crawl result
ENV OUTPUT_DIR=
ENTRYPOINT [ "/bin/grawl" ]
