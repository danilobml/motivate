# Motivate — Quote API

A Go API that serves random motivational quotes.  
It supports seeding quotes from a local JSON file or directly from the ZenQuotes.io API,  
and can optionally send a random quote to a friend by email.

## Overview

Motivate is a small API designed to demonstrate idiomatic Go web development using only the standard library and a few minimal helpers.

Features
- Serve random quotes via /quote
- Add your own quotes via /add
- Send a random quote by E-mail
- Optional seeding:
  - From a local JSON file (--seed-file)
  - From the ZenQuotes API (--seed-api)
- Middleware for logging, panic recovery, CORS, and request IDs
- Unit tests using httptest

## Running the Application

### Prerequisites
- Go 1.21+ (Go 1.24 recommended)
- Internet access (for ZenQuotes seeding)
- Optional: .env file for environment variables (necessary for email)

### Using Make

| Command | Description |
|----------|--------------|
| `make run` | Builds and runs the API normally (no seeding) |
| `make run_seedfile` | Builds, runs and seeds from `./seed_quotes.json` |
| `make run_seedapi` | Builds, runs and seeds from the ZenQuotes API |
| `make test` | Runs all tests under `/test` with verbose output |

## CLI Options

```
go run ./cmd/api [flags]
```

| Flag | Type | Description |
|------|------|-------------|
| `--seed-file` | string | Path to a local JSON file containing quotes |
| `--seed-api` | bool | Fetch quotes from the ZenQuotes.io API |
| *(none)* | | Start empty (no quotes) |

## API Endpoints

| Method | Path | Description |
|--------|------|--------------|
| `GET` | `/health` | Health check (`ok`) |
| `GET` | `/quote` | Returns a random quote (404 if none available) |
| `POST` | `/add` | Add a quote: `{ "text": "...", "author": "..." }` |
| `POST` | `/share` | Send a random quote via email: `{ "to": ["user@example.com"] }` |

### Example: Add a quote
```
curl -X POST http://localhost:8080/add   -H "Content-Type: application/json"   -d '{"text": "Do or do not, there is no try.", "author": "Yoda"}'
```

Response:
```
{
  "id": "e7a5d4fa-bf2d-4e1a-9f1a-4d9e3c8c0f1b",
  "text": "Do or do not, there is no try.",
  "author": "Yoda"
}
```

### Example: Fetch a random quote
```
curl http://localhost:8080/quote
```

Response:
```
{
  "id": "1b2c3d4e-5678-90ab-cdef-1234567890ab",
  "text": "You miss 100% of the shots you don't take.",
  "author": "Wayne Gretzky"
}
```

### Example: Email a random quote
```
curl -X POST http://localhost:8080/share   -H "Content-Type: application/json"   -d '{"to": ["someone@example.com"]}' (can be more than one address: ["someone@example.com", "someone-else@example.com"] )
```

Response (on success): `200 OK`  
If the email service is not configured (no env variables set), returns:
```
{ "error": "mail service disabled" }
```

## Data Seeding

### 1. From Local JSON
Use `--seed-file` or `make run_seedfile`.

Example `seed_quotes.json`:
```
[
  { "text": "Be yourself; everyone else is already taken.", "author": "Oscar Wilde" },
  { "text": "The best revenge is massive success.", "author": "Frank Sinatra" }
]
```

### 2. From ZenQuotes API
Use `--seed-api` or `make run_seedapi`.

Internally:
- `ZenQuoteRepository` fetches from `https://zenquotes.io/api/quotes`
- `ZenQuoteService` converts each into a `Quote` object
- Stored in-memory

## Email Configuration

To enable email delivery, set these environment variables in `.env` or your shell:

| Variable | Description | Example |
|-----------|--------------|---------|
| `FROM_EMAIL` | Sender email address | `motivate@gmail.com` |
| `FROM_EMAIL_PASSWORD` | Password or app-specific password | `app-pass-1234` |
| `FROM_EMAIL_SMTP` | SMTP username (for auth) | `smtp.gmail.com` |
| `SMTP_ADDR` | SMTP host and port | `smtp.gmail.com:587` | (for Gmail)

If any of these are missing, email functionality will be disabled and the `/share` endpoint will return:
```
{ "error": "mail service disabled" }
```
Add your email address to the first variable, get the password (e.g. in Gmail, go to account settings -> app passwords) and add to the second. Modify SMTP and ADDR if not using Gmail.

## Testing

Run all tests:
```
make test
```

## Environment Variables

Example `.env`:
```
PORT=:8080
READ_HEADER_TIMEOUT=5
READ_TIMEOUT=15
WRITE_TIMEOUT=15
IDLE_TIMEOUT=60
FROM_EMAIL=motivate@example.com
FROM_EMAIL_PASSWORD=app-pass-1234
FROM_EMAIL_SMTP=smtp.gmail.com
SMTP_ADDR=smtp.gmail.com:587
```

## Project Tree

```
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
│       ├── mail_service.go
│       ├── quote_service.go
│       └── zenquote_service.go
├── seed_quotes.json
├── test/
│   └── quotes_test.go
├── Makefile
└── README.md
```

## License

MIT License © 2025 Danilo Barolo Martins de Lima
