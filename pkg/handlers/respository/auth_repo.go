package respository

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SignUp(user models.User) error
	UpdateUser(user models.User) error
	GetUserByEmail(email string) (*models.Employee, error)
	GetUserByID(id uint) (*models.Employee, error)
	IncrementRequestCounter(email string, window time.Duration, maxRequests int) (int64, error)
}

type authRepositoryImpl struct {
	tx    *db.TxStore
	mutex sync.Mutex
}

func NewAuthRepositopry(tx *db.TxStore) AuthRepository {
	return &authRepositoryImpl{
		tx: tx,
	}
}

func (r *authRepositoryImpl) GetUserByEmail(email string) (*models.Employee, error) {
	var user models.Employee
	// check email
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		result := db.Where("email = ?", email).First(&user)
		return result.Error
	})

	// match email
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) SignUp(user models.User) error {

	var employee models.Employee = models.Employee{User: user, Position: models.Unknown}
	// create if not exist
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {

		db.AutoMigrate(&models.Employee{})

		return db.Create(&employee).Error
	})
}

func (r *authRepositoryImpl) UpdateUser(user models.User) error {
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		db.AutoMigrate(&models.User{})

		return db.Updates(user).Error
	})
}

func (r *authRepositoryImpl) GetUserByID(id uint) (*models.Employee, error) {
	var user models.Employee
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.First(&user, id).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) IncrementRequestCounter(email string, window time.Duration, maxRequests int) (int64, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var key = fmt.Sprintf("%s:request", email)

	// Increment the request counter atomically
	counter, err := r.tx.RD.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration for the key if not set
	r.tx.RD.Expire(key, window)

	// If the counter exceeds the maximum requests, return an error
	if counter > int64(maxRequests) {
		return counter, fmt.Errorf("too many requests, try again later")
	}

	return counter, nil
}
