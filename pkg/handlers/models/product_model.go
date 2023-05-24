package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

func UnmarshalProduct(data []byte) (Product, error) {
	var r Product
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Product) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Product struct {
	gorm.Model

	ProductName string  `json:"ProductName"`
	Category    string  `json:"Category"`
	Price       float64 `json:"Price"`
	Stock       int     `json:"Stock"`
	StoreID     uint    `gorm:"primaryKey;"`
	ID          uint    `gorm:"primarykey;"`
}
