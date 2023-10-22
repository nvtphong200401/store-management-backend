package respository

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleRepository interface {
	SellItems(sale models.SaleModel) error
	// GetSaleByID(id uint, storeID uint) (int, gin.H)
	GetSaleByID(id uint, storeID uint) (*models.SaleModel, error)
	// GetSales(storeID uint, page, limit int) (int, gin.H)
	GetSales(storeID uint, page, limit int) (models.PaginationModel[models.SaleModel], error)
	BuyItems(sale models.SaleModel) error
}

type saleRepositoryImpl struct {
	tx *db.TxStore
}

func NewSaleRepository(tx *db.TxStore) SaleRepository {
	return &saleRepositoryImpl{
		tx: tx,
	}
}

func (r *saleRepositoryImpl) SellItems(sale models.SaleModel) error {
	// totalPrice := r.calculateTotalPrice(items, CalculatePriceOut)

	// sale := models.SaleModel{
	// 	StoreID:    storeID,
	// 	EmployeeID: employeeID,
	// 	TotalPrice: totalPrice,
	// 	SaleItems:  items,
	// }
	return r.tx.WriteData("", func(db *gorm.DB) (interface{}, error) {
		db.AutoMigrate(&models.SaleModel{})
		db.AutoMigrate(&models.SaleItem{})
		for _, si := range sale.SaleItems {
			var product models.Product
			if e := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, si.ProductID).Error; e != nil {
				return nil, e
			}
			if si.Stock > uint(product.Stock) {
				return nil, fmt.Errorf("invalid stock")
			}
			product.Stock -= si.Stock
			if e := db.Save(&product).Error; e != nil {
				return nil, e
			}
		}
		err := db.Create(&sale).Error
		if err != nil {
			return nil, err
		}
		return sale, nil
	})
}

func (r *saleRepositoryImpl) BuyItems(sale models.SaleModel) error {
	fmt.Println("sale: ", sale)
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		db.AutoMigrate(&models.SaleModel{})
		db.AutoMigrate(&models.SaleItem{})
		for _, si := range sale.SaleItems {
			var product models.Product
			// Try to find the product by ID
			if e := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, si.ProductID).Error; e != nil {
				if errors.Is(e, gorm.ErrRecordNotFound) {
					// If the product doesn't exist, create a new one

					product = si.Product
					now := time.Now()
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
				} else {
					return e
				}
			} else {
				// If the product exists, update the stock
				product.Stock += si.Stock
				if e := db.Save(&product).Error; e != nil {
					return e
				}
			}
		}
		return db.Create(&sale).Error
	})

}

func (r *saleRepositoryImpl) GetSaleByID(id uint, storeID uint) (*models.SaleModel, error) {
	var sale models.SaleModel

	cacheKey := fmt.Sprintf("store:sale:%d:%d", storeID, id)

	err := r.tx.ReadData(cacheKey, &sale, func(db *gorm.DB) (interface{}, error) {
		e := db.Where("id = ? and store_id = ?", id, storeID).Find(&sale).Error
		if e != nil {
			return nil, e
		}
		return sale, nil
	})
	return &sale, err
}

func (r *saleRepositoryImpl) GetSales(storeID uint, page, limit int) (models.PaginationModel[models.SaleModel], error) {
	var pagination = models.PaginationModel[models.SaleModel]{
		Data:       []models.SaleModel{},
		TotalItems: 0,
	}

	cacheKey := fmt.Sprintf("store:sales:%d:%d:%d", storeID, page, limit)

	err := r.tx.ReadData(cacheKey, &pagination, func(db *gorm.DB) (interface{}, error) {
		// Count total items
		db.Model(&models.SaleModel{}).Count(&pagination.TotalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit
		if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Preload("SaleItems").Find(&pagination.Data).Error; err != nil {
			return nil, err
		}
		return pagination, nil
	})
	if err != nil {
		return pagination, err
	}

	// Prepare metadata

	return pagination, nil
}
