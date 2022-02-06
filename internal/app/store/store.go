package store

import (
	"github.com/parMaster/ethtrx/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
)

type Storer interface {
	Create(*model.Transaction) error
	Find(hash string) (*model.Transaction, error)
	BlockExists(number string) bool
	Update(*model.Transaction) error
	MostRecentBlock() (string, error)
	FindTx(bson.M, int64, int64) ([]model.Transaction, error)
}
