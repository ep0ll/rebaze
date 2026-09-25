.PHONY: build test tidy lint clean install ci

BINARY := bin/rebaze
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0-dev")
LDFLAGS := -ldflags "-X github.com/ep0ll/rebaze.Version=$(VERSION) -X github.com/ep0ll/rebaze.Revision=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)"

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BINARY) ./cli

install:
	go install $(LDFLAGS) ./cli

test:
	go test ./... -count=1 -race -timeout 120s

tidy:
	go mod tidy

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed — skipping"

clean:
	rm -rf bin/ coverage.out

ci: tidy build test
