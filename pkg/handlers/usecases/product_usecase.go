package usecases

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
)

type ProductUseCases struct {
	product_repo respository.ProductRepository
}

func NewProductUseCases(product_repo respository.ProductRepository) ProductUseCases {
	return ProductUseCases{
		product_repo: product_repo,
	}
}

func (p *ProductUseCases) GetProducts(storeID uint, page int, limit int) (int, gin.H) {
	totalItems, data, err := p.product_repo.GetProducts(storeID, page, limit)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithPagination(totalItems, limit, page, data)
}

func (p *ProductUseCases) UpdateProducts(products []models.Product) (int, gin.H) {
	for index := range products {
		products[index].UpdatedAt = time.Now()
	}
	err := p.product_repo.UpdateProduct(products)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithResult(products)
}

func (p *ProductUseCases) SearchProducts(keyword string, storeID uint, page int, limit int) (int, gin.H) {
	totalItems, data, err := p.product_repo.SearchProducts(keyword, storeID, page, limit)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithPagination(totalItems, limit, page, data)
}

func (p *ProductUseCases) DeleteProducts(products []models.Product) (int, gin.H) {
	err := p.product_repo.DeleteProduct(products)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithResult("success")
}
