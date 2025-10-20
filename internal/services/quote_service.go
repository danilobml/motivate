package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/danilobml/motivate/internal/errs"
	"github.com/danilobml/motivate/internal/models"
	"github.com/danilobml/motivate/internal/repositories"
)

type QuoteService struct {
	quoteRepository *repositories.InMemoryQuoteRepository
}

func NewQuoteService(repo *repositories.InMemoryQuoteRepository) *QuoteService {
	return &QuoteService{
		quoteRepository: repo,
	}
}

func (qs *QuoteService) GetRandomQuote() (*models.Quote, error) {
	quotes := qs.quoteRepository.List()

	if len(quotes) == 0 {
		return nil, errs.ErrEmpty
	}

	index := rand.Intn(len(quotes))

	return &quotes[index], nil
}

func (qs *QuoteService) CreateQuote(text, author string) (*models.Quote, error) {
	id := uuid.New().String()

	if author == "" {
		author = "Unknown"
	}

	newQuote := models.Quote{
		Id: id,
		Text: text,
		Author: author,
	}

	quote, err := qs.quoteRepository.Save(newQuote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}


func (qs *QuoteService) SeedDbFromFile(filePath string) error {
	start := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		message := fmt.Sprintf("Failed to open json file: %s", err.Error())
		return errors.New(message)
	}
	defer file.Close()

	quotes := []models.Quote{}
	err = json.NewDecoder(file).Decode(&quotes)
	if err != nil {
		return err
	}

	for _, quote := range quotes {
		newQuote := models.Quote{
			Id: quote.Id,
			Text: quote.Text,
			Author: quote.Author,
		}
		qs.quoteRepository.Save(newQuote)
	}

	elapsed := time.Since(start)
	log.Printf("Quotes DB seeded successfully from file! Quotes loaded: %d. Elapsed time: %v.\n", len(quotes), elapsed)

	return nil
}
