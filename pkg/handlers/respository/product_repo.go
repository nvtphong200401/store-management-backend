package respository

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	// Deprecated: Use BuyItems in SaleRepo to add new product
	AddProduct(p []models.Product) error
	GetProducts(storeID uint, page int, limit int) (models.PaginationModel[models.Product], error)
	// [UPDATE] product will update any column except stock
	UpdateProduct(p []models.Product) error
	DeleteProduct(p []models.Product) error
	SearchProducts(keyword string, storeID uint, page int, limit int) (models.PaginationModel[models.Product], error)
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

func (r *productRepositoryImpl) GetProducts(storeID uint, page int, limit int) (models.PaginationModel[models.Product], error) {

	// var totalPages int = 0
	var pagination models.PaginationModel[models.Product] = models.PaginationModel[models.Product]{
		Data:       []models.Product{},
		TotalItems: 0,
	}

	cacheKey := fmt.Sprintf("store:products:%d:%d:%d", storeID, page, limit)

	err := r.tx.ReadData(cacheKey, &pagination, func(db *gorm.DB) (interface{}, error) {
		// Count total items
		db.Model(&models.Product{}).Where("store_id = ?", storeID).Count(&pagination.TotalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit
		if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Order("ID").Find(&pagination.Data).Error; err != nil {
			return pagination, err
		}
		return pagination, nil
	})
	return pagination, err
}

func (r *productRepositoryImpl) UpdateProduct(products []models.Product) error {
	cacheKey := fmt.Sprintf("store:products:%d", products[0].StoreID)
	return r.tx.WriteData(cacheKey, func(db *gorm.DB) (interface{}, error) {
		err := db.Save(&products).Error
		return products, err
	})
}

func (r *productRepositoryImpl) DeleteProduct(p []models.Product) error {
	cacheKey := fmt.Sprintf("store:products:%d", p[0].StoreID)
	return r.tx.WriteData(cacheKey, func(db *gorm.DB) (interface{}, error) {
		for _, prod := range p {

			if err := db.Where("id = ? and store_id = ?", prod.ID, prod.StoreID).Delete(&models.Product{}).Error; err != nil {
				return nil, err
			}
		}
		return []models.Product{}, nil
	})
}

func (r *productRepositoryImpl) SearchProducts(keyword string, storeID uint, page int, limit int) (models.PaginationModel[models.Product], error) {

	cacheKey := fmt.Sprintf("store:products:%v:%d:%d:%d", keyword, storeID, page, limit)
	var pagination models.PaginationModel[models.Product] = models.PaginationModel[models.Product]{
		Data:       []models.Product{},
		TotalItems: 0,
	}
	err := r.tx.ReadData(cacheKey, &pagination, func(db *gorm.DB) (interface{}, error) {
		// Count total items
		db.Model(&models.Product{}).Where("store_id = ? AND to_tsvector('fr', product_name || ' ' || id) @@ to_tsquery('fr', ?)", storeID, keyword).Count(&pagination.TotalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit

		if err := db.Limit(limit).Offset(offset).Where("store_id = ? AND to_tsvector('fr', product_name || ' ' || id) @@ to_tsquery('fr', ?)", storeID, keyword).Find(&pagination.Data).Error; err != nil {
			return pagination, err
		}
		// Calculate total pages
		// totalPages = int(int(totalItems)/limit) + 1
		return pagination, nil
	})
	return pagination, err
}
