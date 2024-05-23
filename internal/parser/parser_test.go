package parser

import (
	"ethereum-parser/shared"
	"testing"
)

// TestGetCurrentBlock verifies that the current block number is retrieved correctly
func TestGetCurrentBlock(t *testing.T) {
	cfg := shared.Config{RPCURL: "https://cloudflare-eth.com"}
	parser := NewParser(cfg)

	blockHex := parser.GetCurrentBlock()
	if blockHex == "0" {
		t.Errorf("expected a non-zero block number, but got %s", blockHex)
	}
}

// TestSubscribe verifies that addresses can be subscribed correctly
func TestSubscribe(t *testing.T) {
	cfg := shared.Config{RPCURL: "https://cloudflare-eth.com"}
	parser := NewParser(cfg)

	address := "0x123"

	subscribed := parser.Subscribe(address)
	if !subscribed {
		t.Errorf("expected Subscribe to return true for a new address")
	}

	// Try subscribing again to the same address
	subscribed = parser.Subscribe(address)
	if subscribed {
		t.Errorf("expected Subscribe to return false for an already subscribed address")
	}
}

// TestGetTransactions verifies that transactions for a subscribed address are retrieved correctly
func TestGetTransactions(t *testing.T) {
	cfg := shared.Config{RPCURL: "https://cloudflare-eth.com"}
	parser := NewParser(cfg)

	address := "0x123"
	parser.subscribed[address] = []Transaction{
		{From: "0x123", To: "0x456", BlockNumber: "0x1", TransactionIndex: "0x1"},
	}

	transactions := parser.GetTransactions(address)
	if len(transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(transactions))
	}
	if transactions[0].From != "0x123" {
		t.Errorf("expected transaction From to be '0x123', got '%s'", transactions[0].From)
	}
	if transactions[0].To != "0x456" {
		t.Errorf("expected transaction To to be '0x456', got '%s'", transactions[0].To)
	}
	if transactions[0].BlockNumber != "0x1" {
		t.Errorf("expected transaction BlockNumber to be '0x1', got '%s'", transactions[0].BlockNumber)
	}
	if transactions[0].TransactionIndex != "0x1" {
		t.Errorf("expected transaction TransactionIndex to be '0x1', got '%s'", transactions[0].TransactionIndex)
	}
}

// TestFetchCurrentBlock verifies that the current block number is fetched correctly from the Ethereum JSON-RPC endpoint
func TestFetchCurrentBlock(t *testing.T) {
	cfg := shared.Config{RPCURL: "https://cloudflare-eth.com"}
	parser := NewParser(cfg)

	blockNumber, err := parser.fetchCurrentBlock()
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if blockNumber == 0 {
		t.Errorf("expected a non-zero block number, got %d", blockNumber)
	}
}
