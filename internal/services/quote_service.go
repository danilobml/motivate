package services

import (
	"math/rand"

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