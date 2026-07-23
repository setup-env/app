.PHONY: build test race vet fmt validate check

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

validate:
	go run ./cmd/setup-env module validate-catalog
	go run ./cmd/setup-env module validate examples/setup-env.yaml

check: fmt vet test validate build
