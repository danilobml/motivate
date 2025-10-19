package httpx

import (
	"log"
	"net/http"

	"github.com/danilobml/motivate/internal/helpers"
)

func NewServer(handler http.Handler) {
	srv := http.Server{
		Addr:              helpers.GetenvString("PORT", ":8080"),
		ReadHeaderTimeout: helpers.GetenvDuration("READ_HEADER_TIMEOUT", 5),
		ReadTimeout:       helpers.GetenvDuration("READ_TIMEOUT", 15),
		WriteTimeout:      helpers.GetenvDuration("WRITE_TIMEOUT", 15),
		IdleTimeout:       helpers.GetenvDuration("IDLE_TIMEOUT", 60),
		Handler:           handler,
	}

	log.Printf("Server listening on port%s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
