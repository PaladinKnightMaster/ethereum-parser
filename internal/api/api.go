package api

import (
	"encoding/json"
	"ethereum-parser/config"
	"ethereum-parser/internal/parser"
	"log"
	"net/http"
)

var p *parser.Parser

// init initializes the parser with configuration settings
func init() {
	cfg := config.LoadConfig()
	p = parser.NewParser(cfg)
}

// GetCurrentBlock handles HTTP requests to retrieve the current Ethereum block number
func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := p.GetCurrentBlock()
	json.NewEncoder(w).Encode(map[string]string{"current_block": block})
}

// Subscribe handles HTTP requests to subscribe an address for transaction notifications
func Subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success := p.Subscribe(req.Address)
	json.NewEncoder(w).Encode(map[string]bool{"subscribed": success})
}

// GetTransactions handles HTTP requests to retrieve transactions for a subscribed address
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	log.Println("GetTranscations()-", address)
	transactions := p.GetTransactions(address)
	json.NewEncoder(w).Encode(transactions)
}
