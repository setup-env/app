.PHONY: build test race vet fmt check

build:
	go build -o bin/setup-env ./cmd/setup-env

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

check: fmt vet test build
