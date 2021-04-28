.PHONY: all

all: clean test
	go build

test:
	go test

clean:
	rm -f input-initial-configuration

crossbuild: all
	GOOS=linux GOARCH=arm GOARM=7 go build

build: all
	go build