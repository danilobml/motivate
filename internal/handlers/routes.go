package handlers

import (
	"net/http"

	"github.com/danilobml/motivate/internal/httpx/middleware"
)

func RegisterRoutes(qr *QuotesRouter) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", getHealth)
	mux.HandleFunc("GET /quote", qr.getRandomQuote)
	mux.HandleFunc("POST /add", qr.createQuote)
	mux.HandleFunc("POST /share", qr.emailRandomQuote)

	return middleware.Cors(middleware.RequestId(middleware.Logger(middleware.Recover(mux))))
}
