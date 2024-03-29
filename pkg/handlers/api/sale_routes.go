package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

type SaleAPI interface {
	CreateSale(c *gin.Context)
	GetSaleByID(c *gin.Context)
	GetSales(c *gin.Context)
}

type saleAPIImpl struct {
	ss usecases.SaleUseCases
}

func NewSaleAPI(sr usecases.SaleUseCases) SaleAPI {
	return &saleAPIImpl{
		ss: sr,
	}
}

func (api *saleAPIImpl) CreateSale(c *gin.Context) {

	var products []models.Product
	if err := c.ShouldBindBodyWith(&products, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := helpers.GetEmployee(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	statusCode, response := api.ss.SellItems(products, employee.ID, employee.StoreID)

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

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	code, response := api.ss.GetSales(employee.StoreID, page, limit)
	c.JSON(code, response)
}
