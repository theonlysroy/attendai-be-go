.PHONY: tidy build run test

tidy:
	@go mod tidy

build:
	@go build -o bin/attendai .

run: build
	@./bin/attendai

test:
	@go test -v ./...
