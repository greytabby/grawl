NAME := Grawl
VERSION := 1.0.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.version=$(VERSION) -X main.revision=$(REVISION)"

export GO111MODULE=on

## Install dependencies
.PHONY: deps
deps:
	go get -v -d

.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...

bin/%: deps
	go build -ldflags $(LDFLAGS) -o $@

.PHONY: build
build: bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/