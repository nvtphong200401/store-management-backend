package respository

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/helpers"
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
	tx          helpers.TxStore
	redisClient *redis.Client
}

func NewProductRepository(gormClient *gorm.DB, redisClient *redis.Client) ProductRepository {
	return &productRepositoryImpl{
		tx:          helpers.NewTXStore(gormClient),
		redisClient: redisClient,
	}
}

func (r *productRepositoryImpl) AddProduct(p *models.Product) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return r.tx.ExecuteTX(func(db *gorm.DB) error {
		db.AutoMigrate(&models.Product{})
		return db.Create(&p).Error
	})
}

func (r *productRepositoryImpl) GetProducts(storeID uint, page int, limit int) (int, gin.H) {
	var products []models.Product
	var totalItems int64 = 0
	var totalPages int = 0

	err := r.tx.ExecuteTX(func(db *gorm.DB) error {
		// Count total items
		db.Model(&models.Product{}).Count(&totalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit
		if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Find(&products).Error; err != nil {
			return err
		}
		// Calculate total pages
		totalPages = int(int(totalItems)/limit) + 1
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"error": err,
		}
	}

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
	return r.tx.ExecuteTX(func(db *gorm.DB) error {
		return db.Save(&product).Error
	})
}

func (r *productRepositoryImpl) DeleteProduct(id uint, storeID uint) error {
	return r.tx.ExecuteTX(func(db *gorm.DB) error {
		return db.Where("id = ? and store_id = ?", id, storeID).Delete(&models.Product{}).Error
	})
}

func (r *productRepositoryImpl) SearchProduct(keyword string, storeID uint) []models.Product {
	var products []models.Product
	keywordLike := "%" + keyword + "%"
	log.Println(keywordLike)
	err := r.tx.ExecuteTX(func(db *gorm.DB) error {
		return db.Where("(product_name LIKE ? OR id = ?) AND store_id = ?", keywordLike, keyword, storeID).Find(&products).Error
	})

	if err != nil {
		log.Println(err.Error())
		return []models.Product{}
	}
	return products
}
