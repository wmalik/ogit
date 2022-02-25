.DEFAULT_GOAL := all

.PHONY: lint
lint:
	go vet ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -count 1 ./...

.PHONY: run
run:
	go run cmd/ogit/main.go

.PHONY: all
all: build lint test
