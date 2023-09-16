package models

import (
	"time"

	"gorm.io/gorm"
)

type SaleModel struct {
	gorm.Model
	StoreID    uint       `json:"StoreID,-"`
	EmployeeID uint       `json:"-"`
	TotalPrice float64    `json:"TotalPrice"`
	SaleItems  []SaleItem `gorm:"foreignKey:SaleID"`
}

type SaleItem struct {
	SaleID uint `gorm:"primarykey" json:"SaleID,-"`
	// Sale      SaleModel      `gorm:"foreignKey:SaleID;" json:"-"`
	ProductID string `gorm:"primarykey" json:"ID,omitempty"` // product ID
	Product   Product
	Stock     uint           `json:"Stock"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
