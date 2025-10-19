package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/danilobml/motivate/internal/services"
)

type QuotesRouter struct {
	quotesService *services.QuoteService
}

type NewQuoteRequest struct {
	Text   string `json:"text" validate:"required,min=1,max=512"`
	Author string `json:"author" validate:"max=128"`
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

func (qr *QuotesRouter) createQuote(w http.ResponseWriter, r *http.Request) {
	quote := NewQuoteRequest{}
	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
	
	validate := validator.New()

    err = validate.Struct(quote)
    if err != nil {
        errors := err.(validator.ValidationErrors)
        http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
        return
    }

	newQuote, err := qr.quotesService.CreateQuote(quote.Text, quote.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newQuote)
}
