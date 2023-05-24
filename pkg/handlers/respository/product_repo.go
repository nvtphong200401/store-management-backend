package respository

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(p *models.Product) error
	GetProducts(storeID uint, page int, limit int) (int, gin.H)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint, storeID uint) error
	SearchProduct(keyword string, storeID uint) []models.Product
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(gormClient *gorm.DB) ProductRepository {
	return &productRepositoryImpl{
		db: gormClient,
	}
}

func (r *productRepositoryImpl) AddProduct(p *models.Product) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	r.db.AutoMigrate(&models.Product{})
	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryImpl) GetProducts(storeID uint, page int, limit int) (int, gin.H) {
	var products []models.Product
	var totalItems int64
	// Count total items
	r.db.Model(&models.Product{}).Count(&totalItems)

	// Retrieve paginated products
	offset := (page - 1) * limit
	if err := r.db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Find(&products).Error; err != nil {
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

func (r *productRepositoryImpl) UpdateProduct(product *models.Product) error {
	product.UpdatedAt = time.Now()
	result := r.db.Save(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *productRepositoryImpl) DeleteProduct(id uint, storeID uint) error {
	return r.db.Where("id = ? and store_id = ?", id, storeID).Delete(&models.Product{}).Error

}

func (r *productRepositoryImpl) SearchProduct(keyword string, storeID uint) []models.Product {
	var products []models.Product
	keywordLike := "%" + keyword + "%"
	log.Println(keywordLike)
	err := r.db.Where("(product_name LIKE ? OR id = ?) AND store_id = ?", keywordLike, keyword, storeID).Find(&products).Error
	if err != nil {
		log.Println(err.Error())
		return []models.Product{}
	}
	return products
}
