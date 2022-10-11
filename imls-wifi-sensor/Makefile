.PHONY: all clean test

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

ifeq ($(shell uname -m),armv7l)
ENVVARS = GOOS=linux GOARCH=arm GOARM=7
else
ifeq ($(OS),Windows_NT)
ENVVARS = GOOS=windows GOARCH=amd64
else
ENVVARS = GOOS=linux
endif
endif

all: session-counter wifi-hardware-search-cli

check:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1 run ./...

deps:
	go mod download

session-counter: deps
	${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/linux-session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-cli

wifi-hardware-search-windows: deps
	${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-windows

test:
	go test -coverprofile all.out -timeout 45m $(shell go list ./... | grep -v windows)

show-test-coverage:
	go tool cover -html=all.out
