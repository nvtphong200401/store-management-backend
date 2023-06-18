package db

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type TxStore struct {
	db *gorm.DB
	rd *redis.Client
}

func NewTXStore(db *gorm.DB, rd *redis.Client) TxStore {
	return TxStore{
		db: db,
		rd: rd,
	}
}

func (ts *TxStore) ExecuteTX(fn func(db *gorm.DB, rd *redis.Client) error) error {
	tx := ts.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := fn(tx, ts.rd); err != nil {
		tx.Rollback()
		return err
	}
	err := tx.Commit().Error
	return err
}
