package api

import (
	"encoding/json"
	"ethereum-parser/config"
	"ethereum-parser/internal/parser"
	"log"
	"net/http"
)

var p *parser.Parser

func init() {
	cfg := config.LoadConfig()
	p = parser.NewParser(cfg)
}

func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := p.GetCurrentBlock()
	json.NewEncoder(w).Encode(map[string]string{"current_block": block})
}

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

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	log.Println("GetTranscations()-", address)
	transactions := p.GetTransactions(address)
	json.NewEncoder(w).Encode(transactions)
}
