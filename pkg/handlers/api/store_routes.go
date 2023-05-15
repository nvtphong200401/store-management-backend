package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/models"
)

var es models.EmployeeService

func CreateStore(c *gin.Context) {
	var store models.StoreModel
	err := c.BindJSON(&store)
	anyEmployee, exist := c.Get("user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}
	if exist {
		employee := anyEmployee.(models.Employee)
		es.CreateStore(&store, &employee)
		c.Set("user", employee)
		c.JSON(http.StatusOK, gin.H{"message": "Created successfully"})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Unauthorized").Error())
	}
}

func JoinStore(c *gin.Context) {
	storeID := c.Param("id")

	id, err := strconv.Atoi(storeID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	employee, exist := c.Get("user")
	if exist {
		err = es.JoinStore(uint(id), employee.(models.Employee))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("Store ID not found").Error())
		return
	}
}

func GetEmployeeInfo(c *gin.Context) {
	employee, exist := c.Get("user")
	if exist {
		c.JSON(http.StatusOK, map[string]interface{}{"result": employee.(models.Employee)})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Something went wrong").Error())
	}
}

func GetStoreInfo(c *gin.Context) {
	employee, exist := c.Get("user")
	if exist {
		store := es.GetStoreInfo(employee.(models.Employee).StoreID)
		if store != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"result": store})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{"message": "Has not joined a store"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Something went wrong").Error())
	}
}
