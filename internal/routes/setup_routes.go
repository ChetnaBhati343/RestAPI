package routes

import (
	"net/http"

	"github.com/chetnabhati343/RESTAPI/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
}
