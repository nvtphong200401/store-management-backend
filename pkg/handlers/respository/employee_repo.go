package respository

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeRepository interface {
	JoinStore(storeID uint, employee models.Employee) (int, gin.H)
	CreateStore(s *models.StoreModel, employee *models.Employee) error
	GetStoreInfo(storeID uint) *models.StoreModel
	GetJoinRequest(storeID uint) (int, gin.H)
	UpdateJoinRequest(storeID uint, employeeID uint, accept bool) (int, gin.H)
	GetStores(page int, limit int) (int, gin.H)
}
type employeeRepositoryImpl struct {
	tx *db.TxStore
}

func NewEmployeeRepository(tx *db.TxStore) EmployeeRepository {
	return &employeeRepositoryImpl{
		tx: tx,
	}
}

func (r *employeeRepositoryImpl) GetStores(page, limit int) (int, gin.H) {
	var stores []models.StoreModel
	var totalItems int64 = 0
	var totalPages int = 0

	e := r.tx.ReadData("store", &stores, func(db *gorm.DB) (interface{}, error) {
		// Count total items
		db.Model(&models.StoreModel{}).Count(&totalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit

		if err := db.Limit(limit).Offset(offset).Find(&stores).Error; err != nil {
			return nil, err
		}
		// Calculate total pages
		totalPages = int(int(totalItems)/limit) + 1
		return stores, nil
	})
	if e != nil {
		return http.StatusInternalServerError, gin.H{"error": e.Error()}
	}
	// Prepare metadata
	metadata := gin.H{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"data":        stores,
	}

	return http.StatusOK, metadata
}

func (r *employeeRepositoryImpl) JoinStore(storeID uint, employee models.Employee) (int, gin.H) {

	// if err :=
	e := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		var request models.JoinRequest
		err := db.Clauses(clause.Locking{Strength: "FOR NO KEY UPDATE"}).Where("Store_ID = ? AND Employee_ID = ?", storeID, employee.ID).Find(&request).Error
		if err != nil {
			return err
		}
		if request.IsPending() {
			return errors.New("already pending")
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
		return errors.New("already in a store")
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
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

	key := fmt.Sprintf("store/%d", storeID)
	err := r.tx.ReadData(key, &store, func(db *gorm.DB) (interface{}, error) {
		if e := db.First(&store, storeID).Error; e != nil {
			return nil, e
		}
		return store, nil
	})
	if err != nil {
		return nil
	}

	return &store
}

func (r *employeeRepositoryImpl) GetJoinRequest(storeID uint) (int, gin.H) {
	var joinRequests []models.JoinRequest
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
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
	e := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
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
