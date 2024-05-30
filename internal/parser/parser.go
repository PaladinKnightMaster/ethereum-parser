package parser

import (
	"bytes"
	"encoding/json"
	"ethereum-parser/shared"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Redeclare shared types
type Transaction = shared.Transaction
type JSONRPCRequest = shared.JSONRPCRequest
type Block = shared.Block

// Parser struct defines the parser state and configuration
type Parser struct {
	currentBlock int64
	subscribed   map[string][]Transaction
	mu           sync.RWMutex
	httpClient   *http.Client
	rpcURL       string
	// Maintain maps to store the processed transaction indexes and block numbers for each subscribed address
	processedTransactions map[string]map[string]struct{}
}

// NewParser creates a new Parser instance and starts the block scanning routine
func NewParser(cfg shared.Config) *Parser {
	p := &Parser{
		currentBlock: 0,
		subscribed:   make(map[string][]Transaction),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		rpcURL:                cfg.RPCURL,
		processedTransactions: make(map[string]map[string]struct{}),
	}

	// Start block scanning in a separate goroutine
	go p.scanBlocks()

	return p
}

// GetCurrentBlock returns the current block number in hexadecimal format
func (p *Parser) GetCurrentBlock() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Fetch the latest block number
	blockNum, err := p.fetchCurrentBlock()
	if err != nil {
		log.Println("Error fetching current block at GetCurrentBlock():", err)
	}

	p.currentBlock = blockNum

	return shared.CurrentBlockToHex(p.currentBlock)
}

// Subscribe adds an address to the list of subscribed addresses for transaction notifications
func (p *Parser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.subscribed[address]; !exists {
		p.subscribed[address] = []Transaction{}
		return true
	}
	return false
}

// GetTransactions returns the list of transactions for a given subscribed address
func (p *Parser) GetTransactions(address string) []Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.subscribed[address]
}

// fetchCurrentBlock retrieves the latest block number from the Ethereum JSON-RPC endpoint
func (p *Parser) fetchCurrentBlock() (int64, error) {

	type Response struct {
		Result string `json:"result"`
	}

	reqBody := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`
	req, err := http.NewRequest("POST", p.rpcURL, strings.NewReader(reqBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseInt(r.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// scanBlocks continuously fetches and processes new blocks for transactions
func (p *Parser) scanBlocks() {
	for {
		p.mu.Lock()
		startBlock := p.currentBlock
		p.mu.Unlock()

		// Fetch the latest block number
		lastBlockNumber, err := p.fetchCurrentBlock()
		if err != nil {
			log.Println("Error fetching current block at scanBlocks(): ", err)
		}

		if startBlock == 0 {
			continue
		}

		// Process each block from the current to the latest block
		for i := startBlock; i <= lastBlockNumber; i++ {
			blockTransactions, err := p.getBlockTransactions(i)
			if err != nil {
				continue
			}
			p.mu.Lock()
			for _, tx := range blockTransactions {
				txKey := tx.BlockNumber + "_" + tx.TransactionIndex
				// Add the transaction to the sender's subscription if not already processed
				if _, exists := p.subscribed[tx.From]; exists {
					if _, exists := p.processedTransactions[tx.From][txKey]; !exists {
						// Append the new transaction to subscribed address for notification
						p.subscribed[tx.From] = append(p.subscribed[tx.From], tx)
						// Add the block number and transaction index to the processed map
						if p.processedTransactions[tx.From] == nil {
							p.processedTransactions[tx.From] = make(map[string]struct{})
						}
						p.processedTransactions[tx.From][txKey] = struct{}{}
					}
				}
				// Add the transaction to the receiver's subscription if not already processed
				if _, exists := p.subscribed[tx.To]; exists {
					if _, exists := p.processedTransactions[tx.To][txKey]; !exists {
						// Append the new transaction to subscribed address for notification
						p.subscribed[tx.To] = append(p.subscribed[tx.To], tx)
						// Add the block number and transaction index to the processed map
						if p.processedTransactions[tx.To] == nil {
							p.processedTransactions[tx.To] = make(map[string]struct{})
						}
						p.processedTransactions[tx.To][txKey] = struct{}{}
					}
				}
			}
			p.mu.Unlock()
		}
		// Sleep before the next block scan
		time.Sleep(10 * time.Second)
	}
}

// getBlockTransactions fetches the transactions for a given block number
func (p *Parser) getBlockTransactions(blockNumber int64) ([]Transaction, error) {
	type Response struct {
		Result json.RawMessage `json:"result"`
	}

	reqBody, _ := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{shared.CurrentBlockToHex(blockNumber), true},
		ID:      83,
	})

	req, err := http.NewRequest("POST", p.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(r.Result, &block)
	if err != nil {
		return nil, err
	}

	return block.Transactions, nil
}
