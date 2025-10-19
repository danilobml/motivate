package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/danilobml/motivate/internal/handlers"
	"github.com/danilobml/motivate/internal/httpx"
	"github.com/danilobml/motivate/internal/repositories"
	"github.com/danilobml/motivate/internal/services"
)

const webPort = ":8080"

func main() {
	seedFilePath := flag.String("seed-file", "", "If a path to a json file is given, the quotes database will be seeded from it.")
	flag.Parse()

	quotesRepo := repositories.NewInMemoryQuoteRepository()
	quotesService := services.NewQuoteService(quotesRepo)
	quotesRouter := handlers.NewQuotesRouter(quotesService)

	if *seedFilePath == "" || filepath.Ext(*seedFilePath) != ".json" {
		log.Println("No valid json seed file path given. The API will initialize unseeded.")
	} else {
		err := quotesRepo.SeedDbFromFile(*seedFilePath)
		if err != nil {
			log.Printf("Error seeding DB: %s. The API will initialize unseeded.", err.Error())
		}
	}

	httpx.NewServer(webPort, handlers.RegisterRoutes(quotesRouter))
}
