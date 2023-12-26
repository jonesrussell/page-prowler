.PHONY: all build install clean fmt lint test docker-build docker-push

USERNAME := jonesrussell
PROJECTNAME := page-prowler
VERSION := v1.0.0
GO = go
GO_LDFLAGS = -ldflags "-s -w"
BINARY_DIR = bin
BINARY_NAME = page-prowler

all: fmt lint test build

build:
	$(GO) build $(GO_LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) main.go

install: build
	install -D -m 0755 $(BINARY_DIR)/$(BINARY_NAME) $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -f $(BINARY_DIR)/$(BINARY_NAME)

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

test:
	$(GO) test -v -race ./... -cover

tidy:
	$(GO) mod tidy

profile:
	$(GO) test -cpuprofile cpu.pprof -memprofile mem.pprof -bench .

docker-build:
	docker build -t $(USERNAME)/$(PROJECTNAME):$(VERSION) .

docker-push:
	docker push $(USERNAME)/$(PROJECTNAME):$(VERSION)
