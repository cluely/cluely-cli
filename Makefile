VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  = -s -w \
           -X github.com/cluely/cli/internal/version.Version=$(VERSION) \
           -X github.com/cluely/cli/internal/version.Commit=$(COMMIT) \
           -X github.com/cluely/cli/internal/version.Date=$(DATE)

.PHONY: build install test lint clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/cluely .

install:
	go install -ldflags "$(LDFLAGS)" .

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/
