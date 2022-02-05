package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/parMaster/ethtrx/internal/app/store/mongostore"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {

	// Setting up
	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)

	client, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer client.Disconnect(context.TODO())
	db := client.Database("eth")
	store := mongostore.NewStore(db)

	srv := newServer(store, *config)

	// Testing client methods
	blockNum, err := srv.getBlockNumber()
	assert.NoError(t, err)
	assert.NotEmpty(t, blockNum)

	assert.NotEqualValues(t, "Max rate limit reached", blockNum)

	blockInfo, err := srv.getBlockByNumber(blockNum)
	assert.NoError(t, err)
	assert.NotEmpty(t, blockInfo)
	assert.GreaterOrEqual(t, len(blockInfo.Transactions), 1)

	for _, v := range blockInfo.Transactions {
		assert.NotEmpty(t, v)
		srv.logger.Logf("DEBUG trx %s", v)
	}
}

func Test_transactions(t *testing.T) {

	// Setting up
	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)

	client, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer client.Disconnect(context.TODO())
	db := client.Database("eth")
	store := mongostore.NewStore(db)

	s := newServer(store, *config)

	// Testing
	testBlock, err := s.getBlockByNumber("0xd7c548") // https://etherscan.io/block/14140744
	assert.NoError(t, err)
	assert.NotNil(t, testBlock)

	if true && len(testBlock.Transactions) > 20 {
		testBlock.Transactions = testBlock.Transactions[:20]
	}
	s.logger.Logf("DEBUG Transactions in this block: %d", len(testBlock.Transactions))

	for _, v := range testBlock.Transactions {
		tx, err := s.getTxByHash(v, true)
		assert.NoError(t, err)
		assert.NotNil(t, tx)

		foundTx, err := s.store.Find(v)
		assert.NoError(t, err)
		assert.NotNil(t, foundTx)

		s.logger.Logf("DEBUG Getting tx: %s \t %v \t %v", v, (err == nil), (foundTx != nil))

		// staying in API limits
		time.Sleep(250 * time.Millisecond)
	}

}

func Test_VariousSpikes(t *testing.T) {

	jStr := `
	{
		"jsonrpc":"2.0",
		"id":1,
		"result":{
		   "baseFeePerGas":"0x5cfe76044",
		   "difficulty":"0x1b4ac252b8a531",
		   "extraData":"0xd883010a06846765746888676f312e31362e36856c696e7578",
		   "timestamp":"0x6110bab2",
		   "totalDifficulty":"0x612789b0aba90e580f8",
		   "transactions":[
			  "0x40330c87750aa1ba1908a787b9a42d0828e53d73100ef61ae8a4d925329587b5",
			  "0x6fa2208790f1154b81fc805dd7565679d8a8cc26112812ba1767e1af44c35dd4",
			  "0xe31d8a1f28d4ba5a794e877d65f83032e3393809686f53fa805383ab5c2d3a3c",
			  "0xa6a83df3ca7b01c5138ec05be48ff52c7293ba60c839daa55613f6f1c41fdace",
			  "0x4e46edeb68a62dde4ed081fae5efffc1fb5f84957b5b3b558cdf2aa5c2621e17",
			  "0x356ee444241ae2bb4ce9f77cdbf98cda9ffd6da244217f55465716300c425e82",
			  "0x1a4ec2019a3f8b1934069fceff431e1370dcc13f7b2561fe0550cc50ab5f4bbc",
			  "0xad7994bc966aed17be5d0b6252babef3f56e0b3f35833e9ac414b45ed80dac93"
		   ],
		   "transactionsRoot":"0xaceb14fcf363e67d6cdcec0d7808091b764b4428f5fd7e25fb18d222898ef779",
		   "uncles":[
			  "0x9e8622c7bf742bdeaf96c700c07151c1203edaf17a38ea8315b658c2e6d873cd"
		   ]
		}
	}
	`
	type Body struct {
		Timestamp    string   `json:"timestamp,omitempty"`
		Transactions []string `json:"transactions"`
	}
	type Block struct {
		Body `json:"result"`
	}
	var cont Block
	json.Unmarshal([]byte(jStr), &cont)
	fmt.Printf("%+v\n", cont)
	fmt.Printf("%+v\n", cont.Transactions)

	fmt.Printf("\n\n\n")

	jStr2 := `
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
	`
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

	var tx TransactionResponse
	json.Unmarshal([]byte(jStr2), &tx)
	fmt.Printf("%+v\n", tx)
	// fmt.Printf("%+v\n", cont.Transactions)

}
