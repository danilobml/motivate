package httpx

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	waitForShutdown(&srv, 5*time.Second)
}

func waitForShutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("\nGracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}

	log.Println("Shutdown complete.")
}
