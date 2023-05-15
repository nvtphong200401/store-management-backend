package helpers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/models"
)

func GetEmployee(c *gin.Context) (models.Employee, error) {
	anyEmployee, existed := c.Get("user")
	if !existed {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return models.Employee{}, errors.New("Unauthorized")
	}
	return anyEmployee.(models.Employee), nil
}
