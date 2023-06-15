package respository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeRepository interface {
	JoinStore(storeID uint, employee models.Employee) (int, gin.H)
	CreateStore(s *models.StoreModel, employee *models.Employee) error
	GetStoreInfo(storeID uint) *models.StoreModel
	GetJoinRequest(storeID uint) (int, gin.H)
	UpdateJoinRequest(storeID uint, employeeID uint, accept bool) (int, gin.H)
}
type employeeRepositoryImpl struct {
	tx          helpers.TxStore
	redisClient *redis.Client
}

func NewEmployeeRepository(gormClient *gorm.DB, redisClient *redis.Client) EmployeeRepository {
	return &employeeRepositoryImpl{
		tx:          helpers.NewTXStore(gormClient),
		redisClient: redisClient,
	}
}

func (r *employeeRepositoryImpl) JoinStore(storeID uint, employee models.Employee) (int, gin.H) {

	// if err :=
	e := r.tx.ExecuteTX(func(db *gorm.DB) error {
		var request models.JoinRequest
		err := db.Clauses(clause.Locking{Strength: "FOR NO KEY UPDATE"}).Where("Store_ID = ? AND Employee_ID = ?", storeID, employee.ID).Find(&request).Error
		if err != nil {
			return err
		}
		if request.IsPending() {
			return errors.New("Already pending")
		}
		request = models.JoinRequest{
			StoreID:    storeID,
			EmployeeID: employee.ID,
			Status:     models.PendingStatus,
		}
		db.AutoMigrate(&request)
		if err := db.Save(request).Error; err != nil {
			return err
		}
		return nil
	})
	if e != nil {
		return http.StatusInternalServerError, gin.H{"error": e.Error()}
	}
	return http.StatusOK, gin.H{"message": "Requested"}
}

func (r *employeeRepositoryImpl) CreateStore(s *models.StoreModel, employee *models.Employee) error {

	if employee.AlreadyInStore() {
		return errors.New("Already in a store")
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	return r.tx.ExecuteTX(func(db *gorm.DB) error {
		db.AutoMigrate(&models.StoreModel{})
		if err := db.Create(&s).Error; err != nil {
			return err
		}
		employee.StoreID = s.ID
		employee.Position = models.Owner
		if err := db.Save(employee).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *employeeRepositoryImpl) GetStoreInfo(storeID uint) *models.StoreModel {
	var store models.StoreModel
	key := fmt.Sprintf("%d", storeID)
	// get data from redis
	result, err := r.redisClient.Get(key).Result()

	if err == redis.Nil {
		// does not exist in redis, get it from postgres
		e := r.tx.ExecuteTX(func(db *gorm.DB) error {
			return db.First(&store, storeID).Error
		})
		if e != nil {
			return nil
		}
		// set data to redis
		data, _ := json.Marshal(store)
		r.redisClient.Set(key, string(data), 0)
	} else if err != nil {
		// some error occured
		log.Println("Some error" + err.Error())

		return nil
	} else {
		// exist in redis
		e := json.Unmarshal([]byte(result), &store)
		if e != nil {
			return nil
		}
	}

	return &store
}

func (r *employeeRepositoryImpl) GetJoinRequest(storeID uint) (int, gin.H) {
	var joinRequests []models.JoinRequest
	err := r.tx.ExecuteTX(func(db *gorm.DB) error {
		return db.Where("store_id = ?", storeID).Preload("Employee").Find(&joinRequests).Error
	})
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}
	}

	return http.StatusOK, gin.H{"results": joinRequests}
}

func (r *employeeRepositoryImpl) UpdateJoinRequest(storeID uint, employeeID uint, accept bool) (int, gin.H) {
	var status models.StatusEnum
	if accept {
		status = models.AcceptedStatus
	} else {
		status = models.DeniedStatus
	}
	e := r.tx.ExecuteTX(func(db *gorm.DB) error {
		if err := db.Model(&models.JoinRequest{}).Where("store_id = ? and employee_id = ? and status = ?", storeID, employeeID, models.PendingStatus).Update("status", status).Error; err != nil {
			return err
		}
		if err := db.Model(&models.Employee{}).Where("id = ?", employeeID).Updates(map[string]interface{}{"store_id": storeID, "position": models.Staff}).Error; err != nil {
			return err
		}
		return nil
	})
	if e != nil {
		return http.StatusInternalServerError, gin.H{"error": e.Error()}
	}
	return http.StatusOK, gin.H{"message": "Updated successfully"}

}
