package models

import (
	"encoding/json"
	"time"

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
}

type ProductService struct {
}

func (ps *ProductService) AddProduct(p *Product) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) GetProducts() ([]Product, error) {
	var products []Product

	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) UpdateProduct(id uint, product *Product) error {

	product.ID = id
	product.UpdatedAt = time.Now()
	result := db.Save(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *ProductService) DeleteProduct(id uint) error {
	return db.Where("id = ?", id).Delete(&Product{}).Error
}
