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

.PHONY: assets
assets:
	go get github.com/jteeuwen/go-bindata/go-bindata; cd assets; go run github.com/jteeuwen/go-bindata/go-bindata -o robot.go -pkg assets robot/...

.PHONY: assets-debug
assets-debug:
	cd assets; go run github.com/jteeuwen/go-bindata/go-bindata -debug -o robot.go -pkg assets robot/...

.PHONY: clean-assets
clean-assets:
	rm assets/robot.go

.PHONY: build
build: assets
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-sitemap

.PHONY: debug
debug:assets
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-sitemap
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-sitemap

.PHONY: test
test:assets
	go test -race -cover ./...

.PHONY: produce
produce:
	HUMAN_LOG=1 go run cmd/producer/main.go

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:assets
	go test -cover -coverpkg=github.com/ONSdigital/dp-sitemap/... -component