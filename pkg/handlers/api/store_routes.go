package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

// var es models.EmployeeService

type EmployeeAPI interface {
	CreateStore(c *gin.Context)
	JoinStore(c *gin.Context)
	GetEmployeeInfo(c *gin.Context)
	GetStoreInfo(c *gin.Context)
	GetJoinRequest(c *gin.Context)
	UpdateJoinRequest(c *gin.Context)
}

type employeeAPIImpl struct {
	sr respository.EmployeeRepository
}

func NewEmployeeAPI(er respository.EmployeeRepository) EmployeeAPI {
	return &employeeAPIImpl{
		sr: er,
	}
}

func (api *employeeAPIImpl) CreateStore(c *gin.Context) {
	var store models.StoreModel
	err := c.BindJSON(&store)
	anyEmployee, exist := c.Get("user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}
	if exist {
		employee := anyEmployee.(models.Employee)
		if employee.AlreadyInStore() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Already in a store"})
			return
		} else {
			api.sr.CreateStore(&store, &employee)
			c.Set("user", employee)
			c.JSON(http.StatusOK, gin.H{"message": "Created successfully"})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Unauthorized").Error())
	}
}

func (api *employeeAPIImpl) JoinStore(c *gin.Context) {
	storeID := c.Param("id")

	id, err := strconv.Atoi(storeID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	if employee.AlreadyInStore() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Already in a store"})
		return
	}

	code, reponse := api.sr.JoinStore(uint(id), employee)
	c.JSON(code, reponse)
}

func (api *employeeAPIImpl) GetEmployeeInfo(c *gin.Context) {
	employee, err := helpers.GetEmployee(c)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"result": employee})
		return
	}
}

func (api *employeeAPIImpl) GetStoreInfo(c *gin.Context) {
	employee, exist := c.Get("user")
	if exist {
		store := api.sr.GetStoreInfo(employee.(models.Employee).StoreID)
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

func (api *employeeAPIImpl) GetJoinRequest(c *gin.Context) {
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}

	code, response := api.sr.GetJoinRequest(employee.StoreID)
	c.JSON(code, response)

}

func (api *employeeAPIImpl) UpdateJoinRequest(c *gin.Context) {
	employee, err := helpers.GetEmployee(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var body struct {
		Accept bool `json:"accept"`
	}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	code, response := api.sr.UpdateJoinRequest(employee.StoreID, uint(id), body.Accept)
	c.JSON(code, response)
}
