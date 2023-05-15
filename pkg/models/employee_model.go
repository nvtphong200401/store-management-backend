package models

import (
	"errors"
	"time"
)

type Employee struct {
	User
	StoreID uint `json:"StoreID,omitempty"`
}

type EmployeeService struct {
}

func (es *EmployeeService) JoinStore(storeID uint, employee Employee) error {
	if employee.StoreID != 0 {
		return errors.New("Already in a store")
	}

	// store does not exist
	var store StoreModel
	result := db.First(&store, storeID)
	if result.Error != nil {
		return result.Error
	}
	employee.StoreID = storeID

	if result := db.Save(employee); result != nil {
		return result.Error
	}
	return nil
}

func (es *EmployeeService) CreateStore(s *StoreModel, employee *Employee) error {

	if employee.StoreID != 0 {
		return errors.New("Already in a store")
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	db.AutoMigrate(&StoreModel{})
	if err := db.Create(&s).Error; err != nil {
		return err
	}
	employee.StoreID = s.ID
	if err := db.Save(employee).Error; err != nil {
		return err
	}

	return nil
}

func (es *EmployeeService) GetStoreInfo(storeID uint) *StoreModel {
	var store StoreModel
	if err := db.First(&store, storeID).Error; err != nil {
		return nil
	}
	return &store
}
