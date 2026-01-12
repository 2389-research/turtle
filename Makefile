# ABOUTME: Build and test automation for turtle project
# ABOUTME: Provides targets for building, testing, and installing the TUI app

.PHONY: build test test-race test-coverage install clean run

build:
	go build -o turtle ./cmd/turtle

run: build
	./turtle

test:
	go test ./... -v

test-race:
	go test -race ./... -v

test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

install:
	go install ./cmd/turtle

clean:
	rm -f turtle coverage.out coverage.html
	go clean

lint:
	golangci-lint run

.DEFAULT_GOAL := build
