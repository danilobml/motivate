package services

import (
	"log"
	"time"

	"github.com/danilobml/motivate/internal/models"
	"github.com/danilobml/motivate/internal/repositories"
	"github.com/google/uuid"
)

type ZenQuoteService struct {
	zenquoteRepository *repositories.ZenQuoteRepository
	quoteRepository *repositories.InMemoryQuoteRepository
}

func NewZenQuoteService(quoteRepo *repositories.InMemoryQuoteRepository, zenRepo *repositories.ZenQuoteRepository) *ZenQuoteService {
	return &ZenQuoteService{
		zenquoteRepository: zenRepo,
		quoteRepository: quoteRepo,
	}
}

func (zs *ZenQuoteService) SeedDbFromApi() error {
	start := time.Now()

	zenQuotes, err := zs.zenquoteRepository.GetZenquotesFromApi()
	if err != nil {
		return err
	}
	
	for _, zenQuote := range zenQuotes {
		quote := models.Quote{
			Id: uuid.New().String(),
			Text: zenQuote.Text,
			Author: zenQuote.Author,
		}

		_, err := zs.quoteRepository.Save(quote)
		if err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Printf("Quotes DB seeded successfully! Quotes loaded: %d. Elapsed time: %v.\n", len(zenQuotes), elapsed)
	
	return nil
}
