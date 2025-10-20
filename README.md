# Motivate — Quote API

A Go API that serves random motivational quotes.  
It supports seeding quotes from a local JSON file or directly from the ZenQuotes.io API.

**Part of a series of mini-projects to get more knowledge of GO features and ecosystem.**

## Overview

Motivate is a small API designed to demonstrate idiomatic Go web development using only the standard library and a few minimal helpers.

Features
- Serve random quotes via /quote
- Add your own quotes via /add
- Optional seeding:
  - From a local JSON file (--seed-file)
  - From the ZenQuotes API (--seed-api)
- Middleware for logging, panic recovery, CORS, and request IDs
- Unit tests using httptest

## Running the Application

### Prerequisites
- Go 1.21+ (Go 1.24 recommended)
- Internet access (for ZenQuotes seeding)
- Optional: .env file for environment variables

### Using Make

The included Makefile simplifies common tasks:

| Command | Description |
|----------|--------------|
| make run | Run the API normally |
| make run_seedfile | Run and seed from ./seed_quotes.json |
| make run_seedapi | Run and seed from the ZenQuotes API |
| make test | Run all tests under /test with verbose output |

## CLI Options

You can also run the binary directly:

go run ./cmd/api [flags]

| Flag | Type | Description |
|------|------|-------------|
| --seed-file | string | Path to a JSON file containing quotes |
| --seed-api | bool | If set, fetches quotes from the ZenQuotes.io API |
| (none) |  | Runs without seeding (empty DB) |

## API Endpoints

| Method | Path | Description |
|--------|------|--------------|
| GET | /health | Returns ok for health checks |
| GET | /quote | Returns a random quote (404 if empty) |
| POST | /add | Adds a new quote ({ "text": "...", "author": "..." }) |

Example: Add a quote

curl -X POST http://localhost:8080/add   -H "Content-Type: application/json"   -d '{"text": "Do or do not, there is no try.", "author": "Yoda"}'

Response:
{
  "id": "e7a5d4fa-bf2d-4e1a-9f1a-4d9e3c8c0f1b",
  "text": "Do or do not, there is no try.",
  "author": "Yoda"
}

Example: Fetch a random quote

curl http://localhost:8080/quote

Response:
{
  "id": "1b2c3d4e-5678-90ab-cdef-1234567890ab",
  "text": "You miss 100% of the shots you don't take.",
  "author": "Wayne Gretzky"
}

## Data Seeding

### 1. From Local JSON
Use --seed-file or make run_seedfile.

Example seed_quotes.json:
[
  { "text": "Be yourself; everyone else is already taken.", "author": "Oscar Wilde" },
  { "text": "The best revenge is massive success.", "author": "Frank Sinatra" }
]

### 2. From ZenQuotes API
Use --seed-api or make run_seedapi.

Internally:
- ZenQuoteRepository fetches from https://zenquotes.io/api/quotes
- ZenQuoteService transforms data into Quote models
- Saved into in-memory repository

## Middleware

| Middleware | Purpose |
|-------------|----------|
| Logger | Logs method, path, status, and latency |
| Recover | Prevents panics from crashing the server |
| RequestID | Attaches a unique X-Request-ID to each request |
| CORS | Allows cross-origin requests (configurable) |

## Testing

Tests use Go’s httptest package.

Run all tests:
make test

## Environment Variables

The project uses github.com/joho/godotenv to load .env if present.  
You can define variables like:

´´´txt
PORT=:8080
READ_HEADER_TIMEOUT=5
READ_TIMEOUT=15
WRITE_TIMEOUT=15
IDLE_TIMEOUT=60
```

(Currently defaults are hardcoded in middleware; environment support is extendable.)

## Project Tree

.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── quotes.go
│   ├── httpx/
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── recover.go
│   │   └── request_id.go
│   ├── models/
│   │   ├── quote.go
│   │   └── zenquote.go
│   ├── repositories/
│   │   ├── quote_repository.go
│   │   └── zenquotes_repository.go
│   └── services/
│       ├── quote_service.go
│       └── zenquote_service.go
├── seed_quotes.json
├── test/
│   └── quotes_test.go
├── Makefile
└── README.md

## License

MIT License © 2025 Danilo Barolo Martins de Lima
