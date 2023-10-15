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
	ID          string  `gorm:"primarykey" json:"ID"`
	ProductName string  `json:"ProductName"`
	Category    string  `json:"Category"`
	PriceIn     float64 `json:"PriceIn"`
	PriceOut    float64 `json:"PriceOut"`
	Stock       uint    `json:"Stock"`
	StoreID     uint    `gorm:"primaryKey;"`
}
