.PHONY: build clean doc fmt lint test vet install bench

install:
	go get -t -v ./...

build: clean \
	lint \
	vet \
	test \
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
