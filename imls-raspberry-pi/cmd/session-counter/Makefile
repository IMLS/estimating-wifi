.PHONY: all test clean

all: clean test
	go build session-counter.go

clean:
	rm -f session-counter

install: all
	echo Copying executable into /usr/local/bin
	sudo cp session-counter /usr/local/bin/session-counter

test:
	go test api/*.go
	go test model/*.go
	go test -timeout 10m *.go

crossbuild: clean
	GOOS=linux GOARCH=arm GOARM=7 go build
