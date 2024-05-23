package shared

// Config defines the configuration parameters for the parser
type Config struct {
	ServerPort string
	RPCURL     string
}

// JSONRPCRequest represents a JSON-RPC request payload
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// Block represents an Ethereum block with its transactions
type Block struct {
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
}
