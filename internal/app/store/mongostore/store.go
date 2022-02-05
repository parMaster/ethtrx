package mongostore

import (
	"context"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store Creates, Finds and Updates Transactions
type Store struct {
	db *mongo.Database
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Create(tx *model.Transaction) error {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	_, err := txs.InsertOne(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Find(hash string) (tx *model.Transaction, err error) {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	result := txs.FindOne(ctx, bson.M{"hash": hash})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	err = result.Decode(&tx)

	return
}

func (s *Store) Update(tx *model.Transaction) error {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	result, err := txs.UpdateOne(ctx, bson.M{"hash": tx.Hash}, tx)

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return store.ErrRecordNotFound
	}

	lgr.Default().Logf("tx: %+v", tx)

	return nil
}
