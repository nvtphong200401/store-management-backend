package respository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type SaleRepository interface {
	PurchaseItems(items []models.SaleItem, employeeID uint, storeID uint) (int, gin.H)
	GetSaleByID(id uint, storeID uint) (int, gin.H)
	GetSales(storeID uint) (int, gin.H)
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

func (r *saleRepositoryImpl) PurchaseItems(items []models.SaleItem, employeeID uint, storeID uint) (int, gin.H) {
	totalPrice := r.calculateTotalPrice(items)
	saleid, err := r.createSale(storeID, employeeID, totalPrice)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}

	err = r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		db.AutoMigrate(&models.SaleItem{})
		for index := range items {
			items[index].SaleID = saleid

		}
		return db.Create(items).Error
	})

	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}

	return http.StatusOK, gin.H{"SaleID": saleid, "Total Price": totalPrice}
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
		total += product.Price * float64(item.Quantity)
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
		return db.Where("Sale_ID = ?", id).Find(&items).Error
	})

	if err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}
	return http.StatusOK, gin.H{
		"items":       items,
		"total_price": sale.TotalPrice,
	}
}

func (r *saleRepositoryImpl) GetSales(storeID uint) (int, gin.H) {
	var sales []models.SaleModel

	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.Where("store_id = ?", storeID).Find(&sales).Error
	})
	if err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}

	return http.StatusOK, gin.H{
		"sales": sales,
	}
}
