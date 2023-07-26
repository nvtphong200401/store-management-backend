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

type ProductAPI interface {
	InsertProduct(c *gin.Context)
	ListProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	SearchProduct(c *gin.Context)
}
type productAPIImpl struct {
	ps respository.ProductRepository
}

func NewProductAPI(ps respository.ProductRepository) ProductAPI {
	return &productAPIImpl{
		ps: ps,
	}
}

func (api *productAPIImpl) InsertProduct(c *gin.Context) {
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
	// product.StoreID = employee.StoreID
	api.ps.AddProduct(products)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, products)
}

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	var product models.Product

	if err = c.ShouldBindBodyWith(&product, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	product.ID = uint(id)
	product.StoreID = employee.StoreID

	if err = api.ps.UpdateProduct(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot update product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (api *productAPIImpl) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}

	err = api.ps.DeleteProduct(uint(id), employee.StoreID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}

func (api *productAPIImpl) SearchProduct(c *gin.Context) {
	var body struct {
		Keyword string `json:"Keyword"`
	}

	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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

	statusCode, metadata := api.ps.SearchProduct(body.Keyword, employee.StoreID, page, limit)
	c.JSON(statusCode, metadata)
}
