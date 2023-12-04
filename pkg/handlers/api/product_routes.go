package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

type ProductAPI interface {
	// InsertProduct(c *gin.Context)
	ListProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	SearchProduct(c *gin.Context)
}
type productAPIImpl struct {
	ps usecases.ProductUseCases
}

func NewProductAPI(ps usecases.ProductUseCases) ProductAPI {
	return &productAPIImpl{
		ps: ps,
	}
}

// func (api *productAPIImpl) InsertProduct(c *gin.Context) {
// 	var products []models.Product
// 	err := c.BindJSON(&products)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	employee, err := helpers.GetEmployee(c)
// 	if err != nil {
// 		return
// 	}
// 	for index := range products {
// 		products[index].StoreID = employee.StoreID
// 	}
// 	// product.StoreID = employee.StoreID
// 	err = api.ps.AddProduct(products)

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, products)
// }

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Security		ApiKeyAuth
// @Success 200 {string} Helloworld
// @Router /products [get]
func (api *productAPIImpl) ListProduct(c *gin.Context) {
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
	statusCode, metadata := api.ps.GetProducts(employee.StoreID, page, limit)

	c.JSON(statusCode, metadata)
}

func (api *productAPIImpl) UpdateProduct(c *gin.Context) {
	var products []models.Product
	err := c.BindJSON(&products)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}

	for index := range products {
		products[index].StoreID = employee.StoreID
	}

	statusCode, response := api.ps.UpdateProducts(products)
	c.JSON(statusCode, response)

	// if err = api.ps.UpdateProducts(products); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot update product"})
	// 	return
	// }
	// c.JSON(http.StatusOK, products)
}

func (api *productAPIImpl) DeleteProduct(c *gin.Context) {
	var products []models.Product
	err := c.BindJSON(&products)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}

	for index := range products {
		products[index].StoreID = employee.StoreID
	}

	statusCode, response := api.ps.DeleteProducts(products)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(statusCode, response)
}

func (api *productAPIImpl) SearchProduct(c *gin.Context) {

	keyword := c.Query("keyword")

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

	statusCode, metadata := api.ps.SearchProducts(keyword, employee.StoreID, page, limit)
	c.JSON(statusCode, metadata)
}
