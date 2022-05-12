VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"
RELEASEPATH = ../release/bin/*

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

check-install:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1

check:
	golangci-lint run ./...

all: input-configuration session-counter wifi-hardware-search-cli log-event cache-to-sqlite waterfall

deps:
	go mod download

input-configuration: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/input-initial-configuration

session-counter: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-cli

log-event: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/log-event

cache-to-sqlite: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/cache-to-sqlite

waterfall: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/waterfall

clean:
	mkdir -p ../release/{bin,pkg}
	chmod -R +w ../release/{bin,pkg}*
	rm -rf ../release/{bin,pkg}

test:
	make -C cmd/cache-to-sqlite/ test
	make -C cmd/input-initial-configuration/ test
	make -C internal/wifi-hardware-search test
	make -C internal/state test
	make -C cmd/session-counter/ test

old_test:
	go test -coverprofile all.out -timeout 45m ./...
