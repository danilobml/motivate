package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danilobml/motivate/internal/services"
)

type QuotesRouter struct {
	quotesService *services.QuoteService
}

func NewQuotesRouter(service *services.QuoteService) *QuotesRouter {
	return &QuotesRouter{
		quotesService: service,
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (qr *QuotesRouter) getRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := qr.quotesService.GetRandomQuote()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}
