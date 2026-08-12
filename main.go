package main

import (
	"net/http"
)

func main() {
	//set up the HTTP server
	mux := http.NewServeMux()

	//server instance
	serverAddr := fmt.Sprintf(";%s", config.ServerPort)
	server := &http.Server{}
}
