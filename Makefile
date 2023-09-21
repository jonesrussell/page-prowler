-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
USERNAME := $(shell whoami)

# Go related variables.
PKG    = github.com/jonesrussell/crawler
PREFIX := /home/$(USERNAME)/.local

# Add the -extldflags "-static" and -tags "netgo" flags to enable static linking.
GO            := GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/bin CGO_ENABLED=0 go
GO_BUILDFLAGS := -tags "netgo"
GO_LDFLAGS    := -s -w -extldflags "-static"

bin/crawler: FORCE
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)/cmd/crawler'

bin/consumer: FORCE
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)/cmd/consumer'

install: FORCE all
	install -D -m 0755 bin/crawler "$(DESTDIR)$(PREFIX)/bin/crawler"
	install -D -m 0755 bin/consumer "$(DESTDIR)$(PREFIX)/bin/consumer"

clean: FORCE
	rm -f -- bin/crawler
	rm -f -- bin/consumer

vendor: FORCE
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: FORCE
