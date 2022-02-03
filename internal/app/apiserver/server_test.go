package apiserver

import (
	"context"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/parMaster/ethtrx/internal/app/store/mongostore"
	"github.com/stretchr/testify/assert"
)

func Test_Client(t *testing.T) {

	config := NewConfig()
	_, err := toml.DecodeFile("../../../configs/apiserver.toml", config)
	assert.NoError(t, err)

	db, err := newDB(config.MongoURI)
	assert.NoError(t, err)

	defer db.Disconnect(context.TODO())
	store := mongostore.NewStore(db)

	srv := newServer(store, *config)

	blockNum, err := srv.getBlockNumber()
	assert.NoError(t, err)
	assert.NotEmpty(t, blockNum)
}
