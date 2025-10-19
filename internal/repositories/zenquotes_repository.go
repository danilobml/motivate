package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/danilobml/motivate/internal/models"
)

type ZenQuoteRepository struct {
	BaseUrl string
}

func NewZenQuoteRepository(baseUrl string) *ZenQuoteRepository {
	return &ZenQuoteRepository{
		BaseUrl: baseUrl,
	}
}

func (zr *ZenQuoteRepository) GetZenquotesFromApi() ([]models.ZenQuote, error) {
	resp, err := http.Get(zr.BaseUrl)
	if err != nil {
		message := fmt.Sprintf("Failed to fetch quotes: %s", err)
		return nil, errors.New(message)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("Failed to read response body: %s", err)
		return nil, errors.New(message)
	}

	var zenQuotes []models.ZenQuote
	err = json.Unmarshal(body, &zenQuotes)
	if err != nil {
		message := fmt.Sprintf("Failed to read response body: %s", err)
		return nil, errors.New(message)
	}

	return zenQuotes, nil
}
