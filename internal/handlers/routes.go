package handlers

import (
	"net/http"

	"github.com/danilobml/motivate/internal/httpx/middleware"
)

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", getHealthHandler)

	return middleware.Logger(middleware.Recover(mux))
}
