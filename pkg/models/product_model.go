package models

import (
	"encoding/json"
	"log"
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
	Stock       int     `json:"Stock"`
	StoreID     uint    `gorm:"primaryKey;"`
	ID          uint    `gorm:"primarykey;uniqueIndex;"`
}

type ProductService struct {
}

func (ps *ProductService) AddProduct(p *Product) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	db.AutoMigrate(&Product{})
	if err := db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) GetProducts(storeID uint) ([]Product, error) {
	var products []Product
	if err := db.Where("store_id = ?", storeID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) UpdateProduct(product *Product) error {
	product.UpdatedAt = time.Now()
	result := db.Save(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *ProductService) DeleteProduct(id uint, storeID uint) error {
	return db.Where("id = ? and store_id = ?", id, storeID).Delete(&Product{}).Error

}

func (ps *ProductService) SearchProduct(keyword string, storeID uint) []Product {
	var products []Product
	keywordLike := "%" + keyword + "%"
	log.Println(keywordLike)
	err := db.Where("(product_name LIKE ? OR 'id' = ?) AND store_id = ?", keywordLike, keyword, storeID).Find(&products).Error
	if err != nil {
		log.Println(err.Error())
		return []Product{}
	}
	return products
}
