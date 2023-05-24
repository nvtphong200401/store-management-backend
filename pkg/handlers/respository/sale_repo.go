package respository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"gorm.io/gorm"
)

type SaleRepository interface {
	PurchaseItems(items []models.SaleItem, employeeID uint, storeID uint) (int, gin.H)
	GetSaleByID(id uint, storeID uint) (int, gin.H)
	GetSales(storeID uint) (int, gin.H)
}

type saleRepositoryImpl struct {
	db gorm.DB
}

func NewSaleRepository(gormClient *gorm.DB) SaleRepository {
	return &saleRepositoryImpl{
		db: *gormClient,
	}
}

func (r *saleRepositoryImpl) createSale(storeID uint, employeeID uint, totalPrice float64) (uint, error) {
	r.db.AutoMigrate(&models.SaleModel{})

	sale := models.SaleModel{StoreID: storeID, EmployeeID: employeeID, TotalPrice: totalPrice}
	if err := r.db.Create(&sale).Error; err != nil {
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
	r.db.AutoMigrate(&models.SaleItem{})
	for index := range items {
		items[index].SaleID = saleid

	}
	if err = r.db.Create(items).Error; err != nil {
		return http.StatusInternalServerError, gin.H{"error": err}
	}

	return http.StatusOK, gin.H{"SaleID": saleid, "Total Price": totalPrice}
}

func (r *saleRepositoryImpl) calculateTotalPrice(items []models.SaleItem) float64 {
	var total float64 = 0
	for _, item := range items {
		var product models.Product
		if err := r.db.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return 0
		}
		total += product.Price * float64(item.Quantity)
	}
	return total
}

func (r *saleRepositoryImpl) GetSaleByID(id uint, storeID uint) (int, gin.H) {
	var sale models.SaleModel
	if err := r.db.Where("id = ? and store_id = ?", id, storeID).Find(&sale).Error; err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}
	var items []models.SaleItem
	if err := r.db.Where("Sale_ID = ?", id).Find(&items).Error; err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}
	return http.StatusOK, gin.H{
		"items":       items,
		"total_price": sale.TotalPrice,
	}
}

func (r *saleRepositoryImpl) GetSales(storeID uint) (int, gin.H) {
	var sales []models.SaleModel

	if err := r.db.Where("store_id = ?", storeID).Find(&sales).Error; err != nil {
		return http.StatusBadRequest, gin.H{"error": err}
	}

	return http.StatusOK, gin.H{
		"sales": sales,
	}
}
