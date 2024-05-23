package main

import (
	"ethereum-parser/config"
	"ethereum-parser/internal/api"
	"log"
	"net/http"
)

// main function initializes the configuration, sets up HTTP routes, and starts the server
func main() {
	// Load configuration settings
	cfg := config.LoadConfig()

	// Define HTTP routes and their handlers
	http.HandleFunc("/current_block", api.GetCurrentBlock)
	http.HandleFunc("/subscribe", api.Subscribe)
	http.HandleFunc("/transactions", api.GetTransactions)

	// Log the server start and listen on the configured port
	log.Printf("Server is running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
