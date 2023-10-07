package db

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
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

func (ts *TxStore) MigrateUp() {
	ts.db.AutoMigrate(&models.Product{}, &models.Employee{}, &models.SaleModel{}, &models.SaleItem{}, &models.JoinRequest{}, models.StoreModel{})
	ts.db.Exec("DROP TEXT SEARCH CONFIGURATION IF EXISTS fr")
	ts.db.Exec("CREATE TEXT SEARCH CONFIGURATION fr ( COPY = french )")
	ts.db.Exec(`ALTER TEXT SEARCH CONFIGURATION fr
	ALTER MAPPING FOR hword, hword_part, word
	WITH unaccent, french_stem;`)
}

func (ts *TxStore) MigrateDown() {
	ts.db.Migrator().DropTable(&models.Product{}, &models.Employee{}, &models.SaleModel{}, &models.SaleItem{}, &models.JoinRequest{}, models.StoreModel{})
}

func (ts *TxStore) CloseStorage() error {
	db, err := ts.db.DB()
	if err != nil {
		return err
	}
	err = ts.rd.Close()
	if err != nil {
		return err
	}
	return db.Close()
}

// Deprecated: Use ReadData or WriteData instead
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

func (ts *TxStore) WriteData(key RedisKey, source interface{}, writeToDB func(db *gorm.DB) error) error {
	tx := ts.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := writeToDB(tx); err != nil {
		tx.Rollback()
		return err
	}
	data, _ := json.Marshal(source)
	ts.rd.Set(string(key), string(data), 3600000)
	err := tx.Commit().Error
	return err
}

func (ts *TxStore) ReadData(key RedisKey, dest interface{}, getFromDB func(db *gorm.DB) error) error {
	result, e := ts.rd.Get(string(key)).Result()

	if e == redis.Nil { // does not exist in redis, get it from postgres
		e := getFromDB(ts.db)

		if e != nil {
			return e
		}
		// set data to redis
		data, _ := json.Marshal(dest)
		ts.rd.Set(string(key), string(data), 3600000)
		return nil
	} else if e != nil {
		// some error occured
		log.Println("Some error" + e.Error())

		return nil
	} else {
		// exist in redis
		e := json.Unmarshal([]byte(result), &dest)
		if e != nil {
			return nil
		}
	}
	return e
}
