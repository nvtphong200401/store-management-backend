package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"github.com/nvtphong200401/store-management/pkg/models"
)

var ps models.ProductService

func InsertProduct(c *gin.Context) {
	var product models.Product
	err := c.BindJSON(&product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	product.StoreID = employee.StoreID
	ps.AddProduct(&product)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func ListProduct(c *gin.Context) {
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	products, err := ps.GetProducts(employee.StoreID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"result": products})
}

func UpdateProduct(c *gin.Context) {
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

	if err = ps.UpdateProduct(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot update product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}

	err = ps.DeleteProduct(uint(id), employee.StoreID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}

func SearchProduct(c *gin.Context) {
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
	result := ps.SearchProduct(body.Keyword, employee.StoreID)
	c.JSON(http.StatusOK, gin.H{"result": result})
}
