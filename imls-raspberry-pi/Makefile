VERSION = $(shell git describe --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

all: input-configuration session-counter wifi-hardware-search-cli

deps:
	go mod download

input-configuration: deps
	GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/input-initial-configuration

session-counter: deps
	GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/session-counter

wifi-hardware-search-cli: deps
	rm -f release/bin/find-ralink
	GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/find-ralink 

test:
	go test ./internal/wifi-hardware-search/lshw/*.go
	go test ./internal/wifi-hardware-search/search/*.go
	go test ./cmd/input-configuration/*.go
	go test ./cmd/session-counting/*.go
	go test ./cmd/session-counting/api/*.go
	go test ./cmd/session-counting/config/*.go
