package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/api"
	"github.com/nvtphong200401/store-management/pkg/handlers/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiv1 := r.Group("api/v1")
	// Authentication
	apiv1.POST("/login", api.Login)
	apiv1.POST("/signup", api.SignUp)

	apiv1.Use(middleware.AuthMiddleware())
	{
		// Employee
		apiv1.GET("user", api.GetEmployeeInfo)
		// Product
		apiv1.POST("/products", api.InsertProduct)
		apiv1.GET("/products", api.ListProduct)
		apiv1.PUT("/products/:id", api.UpdateProduct)
		apiv1.DELETE("/products/:id", api.DeleteProduct)
		// Store
		apiv1.POST("/store", api.CreateStore)
		apiv1.POST("/store/:id", api.JoinStore)
		apiv1.GET("/store", api.GetStoreInfo)
	}

	return r
}
