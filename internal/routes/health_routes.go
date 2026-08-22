package routes

import (
	"net/http"

	"github.com/chetnabhati343/RESTAPI/internal/handlers"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/health", handler.HealthHandler())
}
