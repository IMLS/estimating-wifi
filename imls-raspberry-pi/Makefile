VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

ifeq ($(shell uname -m),armv7l)
RELEASEPATH = release/bin/*
ENVVARS = CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7
else
RELEASEPATH = release/bin/linux_arm/*
ENVVARS = CGO_ENABLED=1 GOOS=linux
endif

all: input-configuration session-counter wifi-hardware-search-cli log-event

deps:
	go mod download

input-configuration: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/input-initial-configuration

session-counter: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/wifi-hardware-search-cli

log-event: deps
	${ENVVARS} GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/log-event

release: all
	git tag ${VERSION}
	git push origin --tags
	gh release create \
		--title '${VERSION} binaries' \
		--notes "Release by $(shell git config user.email) on $(shell date)" \
		${VERSION} ${RELEASEPATH}

# 2020506 MCJ Do we do anything different in dev?
dev: release

clean:
	mkdir -p release/{bin,pkg}
	chmod -R +w release/{bin,pkg}*
	rm -rf release/{bin,pkg}

test:
	go test -coverprofile lshw.out   ./internal/wifi-hardware-search/lshw/*.go
	go test -coverprofile search.out ./internal/wifi-hardware-search/search/*.go
	go test -coverprofile ic.out     ./internal/config/*.go
	go test -coverprofile iic.out    ./cmd/input-configuration/*.go
	go test -coverprofile sca.out    ./cmd/session-counting/api/*.go
	go test -coverprofile sc.out     ./cmd/session-counting/*.go
