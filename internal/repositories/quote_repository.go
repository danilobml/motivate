package repositories

import (
	"errors"
	"slices"

	"github.com/danilobml/motivate/internal/errs"
	"github.com/danilobml/motivate/internal/models"
)

type InMemoryQuoteRepository struct {
	data []models.Quote
}

func NewInMemoryQuoteRepository() *InMemoryQuoteRepository {
	return &InMemoryQuoteRepository{
		data: []models.Quote{
			{
				Id:     "1",
				Text:   "Test quote",
				Author: "TestAuthor",
			},
			{
				Id:     "2",
				Text:   "Test quote 2",
				Author: "TestAuthor 2",
			},
		},
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

func (ir *InMemoryQuoteRepository) Save(q models.Quote) (*models.Quote, error) {
	dbQuote, err := ir.Find(q.Id)

	switch {
	case err == nil:
		dbQuote.Author = q.Author
		dbQuote.Text = q.Text
		return dbQuote, nil
	case errors.Is(err, errs.ErrNotFound):
		ir.data = append(ir.data, q)
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
