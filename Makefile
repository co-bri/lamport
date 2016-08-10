.PHONY: build clean doc fmt lint test vet install bench


all: clean \
	get_pre_req \
	lint \
	vet \
	test \
	install

get_pre_req:
	go get -t -v ./...

install:
	go install

build: 
	go build

clean:
	go clean

doc:
	godoc -http=:6060

lint:
	golint ./...

test:
	go test ./... -cover

bench:
	go test ./... -bench=.

vet:
	go vet ./...
