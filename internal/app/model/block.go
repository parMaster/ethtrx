package model

type Block struct {
	Hash          string        `json:"hash"`
	Timestamp     string        `json:"timestamp"`
	GasUsed       string        `json:"gasUsed"`
	BaseFeePerGas string        `json:"baseFeePerGas"`
	Number        string        `json:"number"`
	Transactions  []Transaction `json:"transactions"`
}
type BlockResponse struct {
	Block `json:"result"`
}
