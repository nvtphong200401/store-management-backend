package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"github.com/nvtphong200401/store-management/pkg/models"
)

var ss models.SaleService

func CreateSale(c *gin.Context) {
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
	if err = ss.PurchaseItems(items, employee.ID, employee.StoreID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
