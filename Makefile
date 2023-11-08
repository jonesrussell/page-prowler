.PHONY: all build install clean fmt lint test docker-build docker-push

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
	golint ./...

test:
	$(GO) test ./... -cover

docker-build:
	docker build -t $(USERNAME)/$(PROJECTNAME):$(VERSION) .

docker-push:
	docker push $(USERNAME)/$(PROJECTNAME):$(VERSION)
