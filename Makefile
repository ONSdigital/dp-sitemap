BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

SITEMAP_CLI_PATH= cmd/sitemapcli

.PHONY: cli
cli:
	go run ${SITEMAP_CLI_PATH}/main.go ${SITEMAP_CLI_PATH}/fake_scroll.go 

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-sitemap

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1
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

.PHONY: generate-sitemap
generate-sitemap:
	go run cmd/sitemapcli/main.go cmd/sitemapcli/fake_scroll.go -robots-file-path=test_robots.txt -scroll-size=10 -scroll-timeout=2000 -api-url=http://localhost -sitemap-index=1 -zebedee-url=http://localhost:8082 -sitemap-file-path=test_sitemap -enable-fake-scroll=true -generate-sitemap=true

.PHONY: update-sitemap
update-sitemap:
	go run cmd/sitemapcli/main.go cmd/sitemapcli/fake_scroll.go -robots-file-path=test_robots.txt -scroll-size=10 -scroll-timeout=2000 -api-url=http://localhost -sitemap-index=1 -zebedee-url=http://localhost:8082 -sitemap-file-path=test_sitemap -enable-fake-scroll=true -update-sitemap=true
