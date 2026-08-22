package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chetnabhati343/RESTAPI/internal/handlers"
	"github.com/chetnabhati343/RESTAPI/internal/routes"
	"github.com/chetnabhati343/RESTAPI/serverconfig"
)

func main() {

	//Load config
	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}
	//set up the HTTP server
	mux := http.NewServeMux()

	///Create a new Handler
	handler := handlers.NewHandlers()

	//Setup Routes
	routes.SetupRoutes(mux, handler)
	//server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}
}
