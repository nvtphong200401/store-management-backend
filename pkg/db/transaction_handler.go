package db

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type TxStore struct {
	db         *gorm.DB
	rd         *redis.Client
	cacheMutex sync.Mutex
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
	keys, err := ts.rd.Keys("*").Result()
	if err != nil {
		log.Println(err)
	}
	// Delete all keys
	for _, key := range keys {
		if err := ts.rd.Del(key).Err(); err != nil {
			log.Printf("Error deleting key %s: %v", key, err)
		}
	}
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

func (ts *TxStore) WriteData(key string, writeToDB func(db *gorm.DB) (interface{}, error)) error {
	tx := ts.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	source, err := writeToDB(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	data, _ := json.Marshal(source)
	if key != "" {
		ts.rd.Set(string(key), string(data), time.Minute*30)
	}
	err = tx.Commit().Error
	return err
}

func (ts *TxStore) ReadData(key string, dest interface{}, getFromDB func(db *gorm.DB) (interface{}, error)) error {
	cacheKey := string(key)
	var err error

	log.Println(reflect.TypeOf(dest).Name())
	err = ts.checkCache(cacheKey, &dest)
	if err == nil {
		return nil
	} else if err != redis.Nil {
		// Some error occurred while reading from Redis
		log.Printf("Error reading from cache: %v", err)
		return err
	}
	log.Println("cache miss")
	// If cache miss, acquire a lock to prevent concurrent database queries
	ts.cacheMutex.Lock()
	defer ts.cacheMutex.Unlock()
	// Recheck the cache after acquiring the lock
	err = ts.checkCache(cacheKey, &dest)
	if err == nil {
		return nil
	}
	// If still not in cache, fetch data from the database

	dest, err = getFromDB(ts.db)

	if err != nil {
		log.Printf("Error fetching data from the database: %v", err)
		return err
	}
	// Store the fetched data in Redis cache with a TTL

	data, err := json.Marshal(dest)
	if err != nil {
		log.Printf("Error marshaling data for cache: %v", err)
		return err
	}

	if err := ts.rd.Set(cacheKey, string(data), time.Minute*30).Err(); err != nil {
		log.Printf("Error storing data in cache: %v", err)
		return err
	}

	return nil
}

func (ts *TxStore) checkCache(key string, dest interface{}) error {

	cachedData, err := ts.rd.Get(key).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &dest); err != nil {
			log.Printf("Cache hit but failed to unmarshal: %v", err)
			return err
		}
		log.Println("hit cache ", dest)
		return nil
	}
	return err
}
