package store

import "github.com/parMaster/ethtrx/internal/app/model"

type Storer interface {
	Create(*model.Transaction) error
	Find(hash string) (*model.Transaction, error)
	Update(*model.Transaction) error
}
