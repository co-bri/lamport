.PHONY: build clean doc fmt lint test vet install bench

install:
	go get -t -v ./...

all: clean \
	lint \
	vet \
	test \
	build

build: 
	go build

clean:
	go clean

doc:
	godoc -http=:6060

lint:
	golint ./...

test:
	go test ./...

bench:
	go test ./... -bench=.

vet:
	go vet ./...
