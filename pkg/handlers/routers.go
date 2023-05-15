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
	apiProduct := apiv1.Group("products")
	apiSale := apiv1.Group("sales")
	// Authentication
	apiv1.POST("/login", api.Login)
	apiv1.POST("/signup", api.SignUp)

	apiv1.Use(middleware.AuthMiddleware())
	{
		// Employee
		apiv1.GET("/user", api.GetEmployeeInfo)
		// Store
		apiv1.POST("/store", api.CreateStore)
		apiv1.POST("/store/:id", api.JoinStore)
		apiv1.GET("/store", api.GetStoreInfo)
	}
	// Product
	apiProduct.Use(middleware.AuthMiddleware(), middleware.StoreMiddleware())
	{
		apiProduct.POST("", api.InsertProduct)
		apiProduct.GET("", api.ListProduct)
		apiProduct.PUT("/:id", api.UpdateProduct)
		apiProduct.DELETE("/:id", api.DeleteProduct)
		apiProduct.GET("/search", api.SearchProduct)
		// apiv1.POST("/products", api.InsertProduct)
		// apiv1.GET("/products", api.ListProduct)
		// apiv1.PUT("/products/:id", api.UpdateProduct)
		// apiv1.DELETE("/products/:id", api.DeleteProduct)
	}
	apiSale.Use(middleware.AuthMiddleware(), middleware.StoreMiddleware())
	{
		apiSale.POST("", api.CreateSale)
	}

	return r
}
