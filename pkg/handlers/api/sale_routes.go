package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

type SaleAPI interface {
	CreateSale(c *gin.Context)
	GetSaleByID(c *gin.Context)
	GetSales(c *gin.Context)
}

type saleAPIImpl struct {
	ss respository.SaleRepository
}

func NewSaleAPI(sr respository.SaleRepository) SaleAPI {
	return &saleAPIImpl{
		ss: sr,
	}
}

func (api *saleAPIImpl) CreateSale(c *gin.Context) {
	var items []models.SaleItem
	if err := c.ShouldBindBodyWith(&items, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := helpers.GetEmployee(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	statusCode, response := api.ss.PurchaseItems(items, employee.ID, employee.StoreID)

	c.JSON(statusCode, response)
}

func (api *saleAPIImpl) GetSaleByID(c *gin.Context) {
	saleid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	code, response := api.ss.GetSaleByID(uint(saleid), employee.StoreID)
	c.JSON(code, response)

}

func (api *saleAPIImpl) GetSales(c *gin.Context) {
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	code, response := api.ss.GetSales(employee.StoreID)
	c.JSON(code, response)
}
