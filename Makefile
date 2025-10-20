.PHONY: test quote run run_seedfile run_seedapi build

BIN := ./bin/motivate

build:
	go build -o $(BIN) ./cmd/api

run: build
	exec $(BIN) ./cmd/api

run_seedfile: build
	exec $(BIN) --seed-file ./seed_quotes.json

run_seedapi: build
	exec $(BIN) --seed-api

test:
	go test ./test -v 
