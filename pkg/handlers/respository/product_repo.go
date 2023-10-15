package respository

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	// Deprecated: Use BuyItems in SaleRepo to add new product
	AddProduct(p []models.Product) error
	GetProducts(storeID uint, page int, limit int) (int, gin.H)
	// [UPDATE] product will update any column except stock
	UpdateProduct(p []models.Product) error
	DeleteProduct(p []models.Product) error
	SearchProduct(keyword string, storeID uint, page int, limit int) (int, gin.H)
}

type productRepositoryImpl struct {
	tx *db.TxStore
}

func NewProductRepository(tx *db.TxStore) ProductRepository {
	return &productRepositoryImpl{
		tx: tx,
	}
}

func (r *productRepositoryImpl) AddProduct(p []models.Product) error {
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		now := time.Now()
		db.AutoMigrate(&models.Product{})
		// First, try to find the existing product by ID, including soft-deleted records

		for _, product := range p {
			product.CreatedAt = now
			product.UpdatedAt = now
			var existingProduct models.Product = product
			if err := db.Unscoped().Where("id = ?", product.ID).FirstOrCreate(&existingProduct).Error; err != nil {
				return err
			}
			if existingProduct.DeletedAt.Valid {
				// If the product exists, update it with the new data
				if err := db.Unscoped().Model(&existingProduct).Save(&product).Error; err != nil {
					return err
				}
			}

		}
		return nil
	})
}

func (r *productRepositoryImpl) GetProducts(storeID uint, page int, limit int) (int, gin.H) {
	var products []models.Product
	var totalItems int64 = 0
	var totalPages int = 0

	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		// Count total items
		db.Model(&models.Product{}).Where("store_id = ?", storeID).Count(&totalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit
		if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Order("ID").Find(&products).Error; err != nil {
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
		"data":        products,
	}

	return http.StatusOK, metadata
}

func (r *productRepositoryImpl) UpdateProduct(products []models.Product) error {

	for index := range products {
		products[index].UpdatedAt = time.Now()
	}
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.Save(products).Error
	})
}

func (r *productRepositoryImpl) DeleteProduct(p []models.Product) error {
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		for _, prod := range p {

			if err := db.Where("id = ? and store_id = ?", prod.ID, prod.StoreID).Delete(&models.Product{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *productRepositoryImpl) SearchProduct(keyword string, storeID uint, page int, limit int) (int, gin.H) {
	var products []models.Product
	var totalItems int64 = 0
	var totalPages int = 0
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		// Count total items
		db.Model(&models.Product{}).Where("store_id = ? AND to_tsvector('fr', product_name || ' ' || id) @@ to_tsquery('fr', ?)", storeID, keyword).Count(&totalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit

		if err := db.Limit(limit).Offset(offset).Where("store_id = ? AND to_tsvector('fr', product_name || ' ' || id) @@ to_tsquery('fr', ?)", storeID, keyword).Find(&products).Error; err != nil {
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
		"data":        products,
	}

	return http.StatusOK, metadata
}
