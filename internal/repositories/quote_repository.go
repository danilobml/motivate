package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/danilobml/motivate/internal/errs"
	"github.com/danilobml/motivate/internal/models"
)

type InMemoryQuoteRepository struct {
	data []models.Quote
}

func NewInMemoryQuoteRepository() *InMemoryQuoteRepository {
	return &InMemoryQuoteRepository{
		data: []models.Quote{},
	}
}

func (ir *InMemoryQuoteRepository) List() []models.Quote {
	return ir.data
}

func (ir *InMemoryQuoteRepository) Find(id string) (*models.Quote, error) {
	for i := range ir.data {
		if ir.data[i].Id == id {
			return &ir.data[i], nil
		}
	}

	return nil, errs.ErrNotFound
}

func (ir *InMemoryQuoteRepository) Save(quote models.Quote) (*models.Quote, error) {
	dbQuote, err := ir.Find(quote.Id)

	switch {
	case err == nil:
		dbQuote.Author = quote.Author
		dbQuote.Text = quote.Text
		return dbQuote, nil
	case errors.Is(err, errs.ErrNotFound):
		ir.data = append(ir.data, quote)
		return &ir.data[len(ir.data)-1], nil
	default:
		return nil, err
	}
}

func (ir *InMemoryQuoteRepository) Delete(id string) error {
	_, err := ir.Find(id)
	if err != nil {
		return err
	}

	ir.data = slices.DeleteFunc(ir.data, func(quote models.Quote) bool {
		return quote.Id == id
	})

	return nil
}

func (ir *InMemoryQuoteRepository) SeedDbFromFile(filePath string) error {
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

	ir.data = append(ir.data, quotes...)

	log.Println("Quotes db seeded successfully!")

	return nil
}
