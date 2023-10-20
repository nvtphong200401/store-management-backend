package usecases

import (
	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/utils"
)

type SaleUseCases struct {
	sale_repo respository.SaleRepository
}

func NewSaleUseCases(saleRepo respository.SaleRepository) SaleUseCases {
	return SaleUseCases{
		sale_repo: saleRepo,
	}
}

func (s *SaleUseCases) SellItems(products []models.Product, storeID uint, employeeID uint) (int, gin.H) {
	items, err := utils.ProductsToSaleItems(products)
	if err != nil {
		return responseWithError(err)
	}
	totalPrice := calculateTotalPrice(items, CalculatePriceOut)

	sale := models.SaleModel{
		StoreID:    storeID,
		EmployeeID: employeeID,
		TotalPrice: totalPrice,
		SaleItems:  items,
	}

	err = s.sale_repo.SellItems(sale)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(sale)
}

func (s *SaleUseCases) BuyItems(products []models.Product, storeID uint, employeeID uint) (int, gin.H) {
	items, err := utils.ProductsToSaleItems(products)
	if err != nil {
		return responseWithResult(err)
	}
	totalPrice := calculateTotalPrice(items, CalculatePriceIn)

	sale := models.SaleModel{
		StoreID:    storeID,
		EmployeeID: employeeID,
		TotalPrice: -totalPrice,
		SaleItems:  items,
	}

	err = s.sale_repo.BuyItems(sale)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithResult(sale)

}

func (s *SaleUseCases) GetSaleByID(id uint, storeID uint) (int, gin.H) {
	sale, err := s.sale_repo.GetSaleByID(id, storeID)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(sale)
}

func (s *SaleUseCases) GetSales(storeID uint, page, limit int) (int, gin.H) {
	totalItems, data, err := s.sale_repo.GetSales(storeID, page, limit)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithPagination(totalItems, limit, page, data)
}

type saleItemCalculator func(item models.SaleItem, product models.Product) float64

func calculateTotalPrice(items []models.SaleItem, calculator saleItemCalculator) float64 {
	var total float64 = 0
	for _, item := range items {
		var product = item.Product
		total += calculator(item, product)
	}
	return total
}

// CalculatePriceOut calculates the price out for a sale item.
func CalculatePriceOut(item models.SaleItem, product models.Product) float64 {
	return product.PriceOut * float64(item.Stock)
}

// CalculatePriceIn calculates the price in for a sale item.
func CalculatePriceIn(item models.SaleItem, product models.Product) float64 {
	// Calculate price in logic here
	return product.PriceIn * float64(item.Stock)
}
