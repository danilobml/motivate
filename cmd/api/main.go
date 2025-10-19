package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/danilobml/motivate/internal/handlers"
	"github.com/danilobml/motivate/internal/httpx"
	"github.com/danilobml/motivate/internal/repositories"
	"github.com/danilobml/motivate/internal/services"
)

func main() {
	godotenv.Load()

	seedFilePath := flag.String("seed-file", "", "Error: No file path provided. Insert the path to a json file containing quotes. The quotes database will be seeded from it.")
	flag.Parse()

	quotesRepo := repositories.NewInMemoryQuoteRepository()
	quotesService := services.NewQuoteService(quotesRepo)
	quotesRouter := handlers.NewQuotesRouter(quotesService)

	if *seedFilePath != "" && filepath.Ext(*seedFilePath) != ".json" {
		log.Println("No valid json seed file path given. The API will initialize unseeded.")
	} else if *seedFilePath != "" {
		err := quotesService.SeedDbFromFile(*seedFilePath)
		if err != nil {
			log.Printf("Error seeding DB: %s. The API will initialize unseeded.", err.Error())
		}
	}

	httpx.NewServer(handlers.RegisterRoutes(quotesRouter))
}
