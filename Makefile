-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
USERNAME := $(shell whoami)

# Go related variables.
PKG = github.com/jonesrussell/crawler
PREFIX = /home/$(USERNAME)/.local

# Add the -extldflags "-static" and -tags "netgo" flags to enable static linking.
GO = CGO_ENABLED=0 go
GO_BUILDFLAGS = -tags "netgo"
GO_LDFLAGS = -s -w -extldflags "-static"

# Build targets.
BINARY_NAMES = crawler consumer
BINARY_DIR = bin
BINARY_FILES = $(addprefix $(BINARY_DIR)/,$(BINARY_NAMES))

.PHONY: all build install clean update-dependencies clean-binaries clean-vendor

all: build

build: $(BINARY_FILES)

$(BINARY_DIR)/%: FORCE
	$(GO) build -o $@ -ldflags '$(GO_LDFLAGS)' '$(PKG)/cmd/$*'

install: build
	install -D -m 0755 $(BINARY_FILES) $(DESTDIR)$(PREFIX)/bin/

clean: clean-binaries clean-vendor

clean-binaries:
	rm -f $(BINARY_FILES)

clean-vendor:
	rm -rf vendor

update-dependencies:
	$(GO) mod tidy
	$(GO) mod vendor

FORCE:
