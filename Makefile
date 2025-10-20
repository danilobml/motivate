.PHONY: test quote run run_seedfile run_seedapi

run:
	go run ./cmd/api

run_seedfile:
	go run ./cmd/api --seed-file ./seed_quotes.json

run_seedapi:
	go run ./cmd/api --seed-api

test:
	go test ./test -v 
