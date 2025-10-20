package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/danilobml/motivate/internal/helpers"
	"github.com/danilobml/motivate/internal/services"
)

type QuotesRouter struct {
	quotesService *services.QuoteService
	mailService   services.Mailer
}

type NewQuoteRequest struct {
	Text   string `json:"text" validate:"required,min=1,max=512"`
	Author string `json:"author" validate:"max=128"`
}

type EmailRequest struct {
	To []string `json:"to" validate:"required,min=1,dive,required,email"`
}

func NewQuotesRouter(service *services.QuoteService, mailService services.Mailer) *QuotesRouter {
	return &QuotesRouter{
		quotesService: service,
		mailService: mailService,
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (qr *QuotesRouter) getRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := qr.quotesService.GetRandomQuote()
	if err != nil {
		helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (qr *QuotesRouter) createQuote(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	quote := NewQuoteRequest{}
	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	validate := validator.New()

	err = validate.Struct(quote)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	text := strings.TrimSpace(quote.Text)
	author := strings.TrimSpace(quote.Author)

	newQuote, err := qr.quotesService.CreateQuote(text, author)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newQuote)
}

func (qr *QuotesRouter) emailRandomQuote(w http.ResponseWriter, r *http.Request) {
	var requestBody EmailRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
	}

	validate := validator.New()

	err = validate.Struct(requestBody)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	quote, err := qr.quotesService.GetRandomQuote()
	if err != nil {
		helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	body := fmt.Sprintf("\"%s\"\n\n - %s", quote.Text, quote.Author)

	err = qr.mailService.SendMail(requestBody.To, "A motivating quote for you", body)
	if err != nil {
		message := fmt.Sprintf("Failed to send email - %s", err.Error())
		helpers.WriteJSONError(w, http.StatusInternalServerError, message)
	}
}
