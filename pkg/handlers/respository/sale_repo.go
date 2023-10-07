package respository

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleRepository interface {
	SellItems(items []models.SaleItem, employeeID uint, storeID uint) (int, gin.H)
	GetSaleByID(id uint, storeID uint) (int, gin.H)
	GetSales(storeID uint, page, limit int) (int, gin.H)
}

type saleRepositoryImpl struct {
	tx *db.TxStore
}

func NewSaleRepository(tx *db.TxStore) SaleRepository {
	return &saleRepositoryImpl{
		tx: tx,
	}
}

func (r *saleRepositoryImpl) createSale(storeID uint, employeeID uint, totalPrice float64) (uint, error) {
	sale := models.SaleModel{StoreID: storeID, EmployeeID: employeeID, TotalPrice: totalPrice}
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		db.AutoMigrate(&models.SaleModel{})
		return db.Create(&sale).Error
	})

	if err != nil {
		return 0, err
	}
	return sale.ID, nil
}

func (r *saleRepositoryImpl) SellItems(items []models.SaleItem, employeeID uint, storeID uint) (int, gin.H) {
	totalPrice := r.calculateTotalPrice(items)

	sale := models.SaleModel{
		StoreID:    storeID,
		EmployeeID: employeeID,
		TotalPrice: totalPrice,
		SaleItems:  items,
	}
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		db.AutoMigrate(&models.SaleModel{})
		db.AutoMigrate(&models.SaleItem{})
		for _, si := range items {
			var product models.Product
			if e := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, si.ProductID).Error; e != nil {
				return e
			}
			if si.Stock > uint(product.Stock) {
				return fmt.Errorf("invalid stock")
			}
			product.Stock -= int(si.Stock)
			if e := db.Save(&product).Error; e != nil {
				return e
			}
		}
		return db.Create(&sale).Error
	})

	// saleid, err := r.createSale(storeID, employeeID, totalPrice)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}

	// err = r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
	// 	db.AutoMigrate(&models.SaleItem{})
	// 	for index := range items {
	// 		items[index].SaleID = saleid

	// 	}
	// 	return db.Create(items).Error
	// })

	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}

	return http.StatusOK, gin.H{"ID": sale.ID, "TotalPrice": totalPrice}
}

func (r *saleRepositoryImpl) calculateTotalPrice(items []models.SaleItem) float64 {
	var total float64 = 0
	for _, item := range items {
		var product models.Product
		err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
			return db.Where("id = ?", item.ProductID).First(&product).Error
		})
		if err != nil {
			return 0
		}
		total += product.PriceOut * float64(item.Stock)
	}
	return total
}

func (r *saleRepositoryImpl) GetSaleByID(id uint, storeID uint) (int, gin.H) {
	var sale models.SaleModel

	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.Where("id = ? and store_id = ?", id, storeID).Find(&sale).Error
	})

	if err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}

	var items []models.SaleItem
	err = r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		if err = db.Where("Sale_ID = ?", id).Preload("Product").Find(&items).Error; err != nil {
			return err
		}
		for index, value := range items {
			var product models.Product
			if err = db.Where("ID = ?", value.ProductID).First(&product).Error; err != nil {
				return err
			}
			items[index].Product = product
		}
		return nil
	})

	if err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}
	return http.StatusOK, gin.H{
		"Items":      items,
		"TotalPrice": sale.TotalPrice,
	}
}

func (r *saleRepositoryImpl) GetSales(storeID uint, page, limit int) (int, gin.H) {
	var sales []models.SaleModel = make([]models.SaleModel, 0)
	var totalItems int64 = 0
	var totalPages int = 1

	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		// Count total items
		db.Model(&models.SaleModel{}).Count(&totalItems)
		// Retrieve paginated products
		offset := (page - 1) * limit
		if err := db.Limit(limit).Offset(offset).Where("store_id = ?", storeID).Find(&sales).Error; err != nil {
			return err
		}
		// Calculate total pages
		totalPages = int(int(totalItems)/limit) + 1
		return nil
	})
	metadata := gin.H{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"data":        sales,
	}
	if err != nil {
		return http.StatusBadRequest, metadata
	}

	// Prepare metadata

	return http.StatusOK, metadata
}
