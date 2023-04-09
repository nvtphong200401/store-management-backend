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
	// Product
	apiv1.Use(middleware.AuthMiddleware())
	{
		apiv1.POST("/products", api.InsertProduct)
		apiv1.GET("/products", api.ListProduct)
		apiv1.PUT("/products/:id", api.UpdateProduct)
		apiv1.DELETE("/products/:id", api.DeleteProduct)
	}

	return r
}
