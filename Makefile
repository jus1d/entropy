.PHONY: build run test lint clean gen

build:
	go build -o .build/api ./cmd/api

run:
	go run ./cmd/api

test: gen
	go test ./...

lint: gen
	golangci-lint run ./...

gen:
	go generate ./...

clean:
	rm -rf .build/
