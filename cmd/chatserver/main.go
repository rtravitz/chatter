package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rtravitz/chatter/cmd/chatserver/handlers"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "localhost:5000"
	}

	return ":" + port
}

func main() {
	log.Println("Starting on port 5000")
	api := handlers.API()
	log.Fatal(http.ListenAndServe(getPort(), api))
}
