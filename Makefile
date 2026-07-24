.PHONY: build test race vet fmt validate release-snapshot release-verify check

RELEASE_VERSION ?= v0.1.0-snapshot

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
	go run ./cmd/setup-env status
	go run ./cmd/setup-env status --json
	go run ./cmd/setup-env

release-snapshot:
	go run ./cmd/release -version $(RELEASE_VERSION) -output dist -clean-owned

release-verify:
	go run ./cmd/release -version $(RELEASE_VERSION) -output dist -verify-only

check: fmt vet test validate build
