package respository

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	JoinStore(storeID uint, employee models.Employee) (int, gin.H)
	CreateStore(s *models.StoreModel, employee *models.Employee) error
	GetStoreInfo(storeID uint) *models.StoreModel
	GetJoinRequest(storeID uint) (int, gin.H)
	UpdateJoinRequest(storeID uint, employeeID uint, accept bool) (int, gin.H)
}
type employeeRepositoryImpl struct {
	db *gorm.DB
}

func NewEmployeeRepository(gormClient *gorm.DB) EmployeeRepository {
	return &employeeRepositoryImpl{
		db: gormClient,
	}
}

func (r *employeeRepositoryImpl) JoinStore(storeID uint, employee models.Employee) (int, gin.H) {

	var request models.JoinRequest
	// if err :=
	r.db.Where("Store_ID = ? AND Employee_ID = ?", storeID, employee.ID).Find(&request)
	// .Error; err != nil {
	// 	return err
	// }
	if request.IsPending() {

		return http.StatusAlreadyReported, gin.H{"message": "Already pending"}
	}
	request = models.JoinRequest{
		StoreID:    storeID,
		EmployeeID: employee.ID,
		Status:     models.PendingStatus,
	}
	r.db.AutoMigrate(&request)

	if err := r.db.Save(request).Error; err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
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
	r.db.AutoMigrate(&models.StoreModel{})
	if err := r.db.Create(&s).Error; err != nil {
		return err
	}
	employee.StoreID = s.ID
	employee.Position = models.Owner
	if err := r.db.Save(employee).Error; err != nil {
		return err
	}

	return nil
}

func (r *employeeRepositoryImpl) GetStoreInfo(storeID uint) *models.StoreModel {
	var store models.StoreModel
	if err := r.db.First(&store, storeID).Error; err != nil {
		return nil
	}
	return &store
}

func (r *employeeRepositoryImpl) GetJoinRequest(storeID uint) (int, gin.H) {
	var joinRequests []models.JoinRequest

	r.db.Where("store_id = ?", storeID).Preload("Employee").Find(&joinRequests)

	return http.StatusOK, gin.H{"results": joinRequests}
}

func (r *employeeRepositoryImpl) UpdateJoinRequest(storeID uint, employeeID uint, accept bool) (int, gin.H) {
	var status models.StatusEnum
	if accept {
		status = models.AcceptedStatus
	} else {
		status = models.DeniedStatus
	}
	if err := r.db.Model(&models.JoinRequest{}).Where("store_id = ? and employee_id = ? and status = ?", storeID, employeeID, models.PendingStatus).Update("status", status).Error; err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}
	if err := r.db.Model(&models.Employee{}).Where("id = ?", employeeID).Updates(map[string]interface{}{"store_id": storeID, "position": models.Staff}).Error; err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}
	return http.StatusOK, gin.H{"message": "Updated successfully"}

}
