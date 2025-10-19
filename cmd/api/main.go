package main

import (
	"github.com/danilobml/motivate/internal/handlers"
	"github.com/danilobml/motivate/internal/httpx"
	"github.com/danilobml/motivate/internal/repositories"
	"github.com/danilobml/motivate/internal/services"
)

const webPort = ":8080"

func main() {

	quotesRepo := repositories.NewInMemoryQuoteRepository()
	quotesService := services.NewQuoteService(quotesRepo)
	quotesRouter := handlers.NewQuotesRouter(quotesService)

	httpx.NewServer(webPort, handlers.RegisterRoutes(quotesRouter))
}
