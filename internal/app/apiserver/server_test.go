package apiserver

import (
	"context"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/parMaster/ethtrx/internal/app/store/mongostore"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {

	// Setting up
	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)

	db, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer db.Disconnect(context.TODO())
	store := mongostore.NewStore(db)

	srv := newServer(store, *config)

	// Testing client methods
	blockNum, err := srv.getBlockNumber()
	assert.NoError(t, err)
	assert.NotEmpty(t, blockNum)

	assert.NotEqualValues(t, "Max rate limit reached", blockNum.Result)

	// eth_getBlockByNumber
	// https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x10d4f&boolean=true&apikey=YourApiKeyToken

	blockInfo, err := srv.getBlockByNumber(blockNum.Result)
	assert.NoError(t, err)
	assert.NotEmpty(t, blockInfo)
	assert.GreaterOrEqual(t, 1, len(blockInfo.Result.Transactions))

	for _, v := range blockInfo.Result.Transactions {
		assert.NotEmpty(t, v)
	}
}
