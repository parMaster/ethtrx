package apiserver

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/parMaster/ethtrx/internal/app/store/mongostore"
	"github.com/stretchr/testify/assert"
)

// Use suite!
// type UnitTestSuite struct {
// 	suite.Suite
// }

// func (s *UnitTestSuite) SetupTest() {
// }

func Test_Client(t *testing.T) {

	// Setting up
	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)

	config.DaemonMode = false //no update goroutine for test mode

	client, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer client.Disconnect(context.TODO())
	db := client.Database("eth")
	store := mongostore.NewStore(db)

	srv := newServer(store, *config)

	// Testing client methods
	blockNum, err := srv.getCurrentBlockNumber()
	assert.NoError(t, err)
	assert.NotEmpty(t, blockNum)

	assert.NotEqualValues(t, "Max rate limit reached", blockNum)

	blockInfo, err := srv.getBlockByNumber("0xd80eb3", true)
	assert.NoError(t, err)
	assert.NotEmpty(t, blockInfo)
	assert.GreaterOrEqual(t, len(blockInfo.Transactions), 1)

	for _, v := range blockInfo.Transactions {
		assert.NotEmpty(t, v)
		// srv.logger.Logf("DEBUG trx %+v", v)
	}
}

func Test_transactions(t *testing.T) {

	// Setting up
	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)
	config.DaemonMode = false //don't start updating goroutine for test mode

	client, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer client.Disconnect(context.TODO())
	db := client.Database("eth")
	store := mongostore.NewStore(db)

	s := newServer(store, *config)

	// Testing
	testBlock, err := s.getBlockByNumber("0xd80eb3", true) // https://etherscan.io/block/14140744
	assert.NoError(t, err)
	assert.NotNil(t, testBlock)

	assert.Equal(t, 9, len(testBlock.Transactions))
	// s.logger.Logf("DEBUG Transactions in this block: %d", len(testBlock.Transactions))

	for _, v := range testBlock.Transactions {
		tx, err := s.getTxByHash(v.Hash, true)
		assert.NoError(t, err)
		assert.NotNil(t, tx)

		foundTx, err := s.store.Find(v.Hash)
		assert.NoError(t, err)
		assert.NotNil(t, foundTx)

		// s.logger.Logf("DEBUG Getting tx: %+v \t %v \t %v", v, (err == nil), (foundTx != nil))

		// staying within API limits
		time.Sleep(200 * time.Millisecond)
	}

}

// Learning how to work with big hex and show those like floats
func Test_Conversion(t *testing.T) {

	v, err := hexutil.DecodeUint64("0x0")
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), v)

	bigHex := "0x56bb4b75481fde000"

	v, err = hexutil.DecodeUint64(bigHex) // too big to DecodeUint64
	assert.Error(t, err)
	assert.Equal(t, uint64(0xffffffffffffffff), v)

	vBig, err := hexutil.DecodeBig(bigHex)
	assert.NoError(t, err)
	assert.NotEmpty(t, vBig)

	vBigFloat := new(big.Float).SetInt(vBig)
	assert.NotEmpty(t, vBigFloat)

	// all together
	vTogether := new(big.Float).SetInt(hexutil.MustDecodeBig(bigHex))
	assert.NotEmpty(t, vTogether)

	conversionFactor := new(big.Float).SetUint64(1000000000000000000)

	bigFloat := vTogether.Quo(vTogether, conversionFactor)
	fmt.Sprintf("%2.8f", bigFloat)

	conv := BigHexToStr()

	fmt.Println(conv(bigHex))
}
