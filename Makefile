BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-sitemap

.PHONY: build-cli
build-cli:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-sitemap-cli cmd/cli-tool/*.go

.PHONY: build-cli-remote
build-cli-remote:
	GOOS=linux GOARCH=amd64 go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-sitemap-cli-remote cmd/cli-tool/*.go

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	golangci-lint run ./...

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-sitemap
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-sitemap

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: produce
produce:
	HUMAN_LOG=1 go run cmd/producer/main.go

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/ONSdigital/dp-sitemap/... -component
