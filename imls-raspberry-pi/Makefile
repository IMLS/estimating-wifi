VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"
ENVVARS = CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7

all: input-configuration session-counter wifi-hardware-search-cli

deps: versioning
	go mod download

input-configuration: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/input-initial-configuration

session-counter: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/find-ralink

release: all
	git tag ${VERSION}
	git push origin --tags
	gh release create \
		--title '${VERSION} binaries' \
		--notes "Release by $(shell git config user.email) on $(shell date)" \
		${VERSION} release/bin/linux_arm/*

clean:
	mkdir -p release/bin release/pkg
	chmod -R +w release/{bin,pkg}*
	rm -rf release/{bin,pkg}

test:
	go test ./internal/wifi-hardware-search/lshw/*.go
	go test ./internal/wifi-hardware-search/search/*.go
	go test ./cmd/input-configuration/*.go
	go test ./cmd/session-counting/*.go
	go test ./cmd/session-counting/api/*.go
	go test ./cmd/session-counting/config/*.go
