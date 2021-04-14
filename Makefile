.PHONY: all

all: clean test
	go build

test:
	go test

clean:
	rm -f input-initial-configuration