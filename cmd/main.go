package main

import (
	"ethereum-parser/config"
	"ethereum-parser/internal/api"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	http.HandleFunc("/current_block", api.GetCurrentBlock)
	http.HandleFunc("/subscribe", api.Subscribe)
	http.HandleFunc("/transactions", api.GetTransactions)

	log.Printf("Server is running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
