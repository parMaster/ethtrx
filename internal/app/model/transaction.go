package model

type Transaction struct {
	Hash        string `json:"hash" bson:"hash"`
	BlockHash   string `json:"blockHash" bson:"blockHash"`
	BlockNumber string `json:"blockNumber" bson:"blockNumber"`
	From        string `json:"from" bson:"from"`
	To          string `json:"to" bson:"to"`
	Value       string `json:"value" bson:"value"`
	Gas         string `json:"gas" bson:"gas"`
	GasPrice    string `json:"gasPrice" bson:"gasPrice"`
	Type        string `json:"type" bson:"type"`
	ChainId     string `json:"chainId" bson:"chainId"`
}
type TransactionResponse struct {
	Transaction `json:"result"`
}

/*
{
	"jsonrpc":"2.0",
	"id":1,
	"result":{
	   "blockHash":"0xf850331061196b8f2b67e1f43aaa9e69504c059d3d3fb9547b04f9ed4d141ab7",
	   "blockNumber":"0xcf2420",
	   "from":"0x00192fb10df37c9fb26829eb2cc623cd1bf599e8",
	   "gas":"0x5208",
	   "gasPrice":"0x19f017ef49",
	   "maxFeePerGas":"0x1f6ea08600",
	   "maxPriorityFeePerGas":"0x3b9aca00",
	   "hash":"0xbc78ab8a9e9a0bca7d0321a27b2c03addeae08ba81ea98b03cd3dd237eabed44",
	   "input":"0x",
	   "nonce":"0x33b79d",
	   "to":"0xc67f4e626ee4d3f272c2fb31bad60761ab55ed9f",
	   "transactionIndex":"0x5b",
	   "value":"0x19755d4ce12c00",
	   "type":"0x2",
	   "accessList":[

	   ],
	   "chainId":"0x1",
	   "v":"0x0",
	   "r":"0xa681faea68ff81d191169010888bbbe90ec3eb903e31b0572cd34f13dae281b9",
	   "s":"0x3f59b0fa5ce6cf38aff2cfeb68e7a503ceda2a72b4442c7e2844d63544383e3"
	}
}
*/
