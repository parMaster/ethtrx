package model

type Transaction struct {
	Hash          string `json:"hash" bson:"hash"`
	BlockHash     string `json:"blockHash" bson:"blockHash"`
	BlockNumber   string `json:"blockNumber" bson:"blockNumber"`
	BlockHeight   uint64 `json:"blockHeight" bson:"blockHeight"`
	BlockTime     string `json:"blockTime,omitempty" bson:"blockTime"`
	From          string `json:"from" bson:"from"`
	To            string `json:"to" bson:"to"`
	Value         string `json:"value" bson:"value"`
	ValueNumber   string `json:"valueNumber" bson:"valueNumber,omitempty"`
	Confirmations uint64 `json:"confirmations" bson:"confirmations,omitempty"`
	Type          string `json:"type" bson:"type"`
	ChainId       string `json:"chainId" bson:"chainId"`
}
type TransactionResponse struct {
	Transaction `json:"result"`
}
