package main

import (
	"github.com/danilobml/motivate/internal/handlers"
	"github.com/danilobml/motivate/internal/httpx"
)

const webPort = ":8080"

func main() {

	httpx.NewServer(webPort, handlers.RegisterRoutes())
}
