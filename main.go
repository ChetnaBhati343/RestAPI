package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//set up the HTTP server
	mux := http.NewServeMux()

	//server instance
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: nil,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed %v", err)
	}
}
