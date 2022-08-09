.PHONY: all clean test

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

ifeq ($(shell uname -m),armv7l)
ENVVARS = CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7
else
ifeq ($(OS),Windows_NT)
# choco install mingw (for now, since we need CGO)
ENVVARS = CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/c/ProgramData/chocolatey/lib/mingw/tools/install/mingw64/bin/x86_64-w64-mingw32-gcc
else
ENVVARS = CGO_ENABLED=1 GOOS=linux
endif
endif

all: session-counter wifi-hardware-search-cli

check:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1 run ./...

deps:
	go mod download

session-counter: deps
	${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-cli

test:
	make -C internal/wifi-hardware-search test
	make -C internal/state test
	make -C cmd/session-counter/ test

old_test:
	go test -coverprofile all.out -timeout 45m ./...
