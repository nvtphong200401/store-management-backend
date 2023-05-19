package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (ps *ProductService) GetProducts(storeID uint, page int, limit int) (int, gin.H) {
	var products []Product
	var totalItems int64
	// Count total items
	db.Model(&Product{}).Count(&totalItems)

	// Retrieve paginated products
	offset := (page - 1) * limit
	if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Find(&products).Error; err != nil {
		return http.StatusInternalServerError, gin.H{
			"error": err,
		}
	}

	// Calculate total pages
	totalPages := int(int(totalItems)/limit) + 1

	// Prepare metadata
	metadata := gin.H{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"products":    products,
	}

	return http.StatusOK, metadata
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
	err := db.Where("(product_name LIKE ? OR id = ?) AND store_id = ?", keywordLike, keyword, storeID).Find(&products).Error
	if err != nil {
		log.Println(err.Error())
		return []Product{}
	}
	return products
}
