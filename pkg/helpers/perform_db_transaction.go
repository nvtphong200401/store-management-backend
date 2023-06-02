package helpers

import (
	"gorm.io/gorm"
)

type TxStore struct {
	db *gorm.DB
}

func NewTXStore(db *gorm.DB) TxStore {
	return TxStore{
		db: db,
	}
}

func (ts *TxStore) ExecuteTX(fn func(db *gorm.DB) error) error {
	tx := ts.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	err := tx.Commit().Error
	return err
}
