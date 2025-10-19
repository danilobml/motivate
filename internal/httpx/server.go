package httpx

import (
	"log"
	"net/http"
)

func NewServer(addr string, handler http.Handler) {
	srv := http.Server{
		Addr: addr,
		Handler: handler,
	}

	log.Printf("Server listening on port%s", addr)

	if err := srv.ListenAndServe(); err != nil {
		panic("Failed to start server")
	}
}
