package models

import "gorm.io/gorm"

type SaleModel struct {
	gorm.Model
	StoreID    uint       `json:"StoreID,-"`
	Store      StoreModel `gorm:"foreignKey:StoreID;" json:"-"`
	EmployeeID uint       `json:"-"`
	Employee   Employee   `gorm:"foreignKey:EmployeeID" json:"-"`
	TotalPrice float64    `json:"TotalPrice"`
}

type SaleItem struct {
	gorm.Model
	SaleID    uint      `json:"SaleID,-"`
	Sale      SaleModel `gorm:"foreignKey:SaleID"`
	ProductID uint      `json:"ProductID,omitempty"`
	Quantity  uint      `json:"Quantity"`
}
