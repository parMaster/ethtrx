package mongostore

import (
	"context"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/parMaster/ethtrx/internal/app/model"
	"github.com/parMaster/ethtrx/internal/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *Store) BlockExists(number string) bool {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	result := txs.FindOne(ctx, bson.M{"blockNumber": number})
	if result.Err() == mongo.ErrNoDocuments {
		return false
	}

	return true
}

func (s *Store) MostRecentBlock() (string, error) {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"blockNumber", -1}})
	findOptions.SetLimit(1)

	cursor, err := txs.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		defer cursor.Close(ctx)
		return "", err
	}

	var result *model.Transaction

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return "", err
		}
		return result.BlockNumber, nil
	}

	return "", store.ErrRecordNotFound
}

func (s *Store) FindTx(filter bson.M, page, limit int64) ([]model.Transaction, error) {
	txs := s.db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// dirty and lazy pagination
	findOptions := options.Find()
	if limit >= 0 && page > 0 {

		if page == 1 {
			findOptions.SetSkip(0)
			findOptions.SetLimit(limit)
		}

		findOptions.SetSkip((page - 1) * limit)
		findOptions.SetLimit(limit)
	}

	cursor, err := txs.Find(ctx, filter, findOptions)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	var result []model.Transaction
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
