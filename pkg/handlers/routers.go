package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/controller"
)

func InitRouter(c controller.AppController) *gin.Engine {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiv1 := r.Group("api/v1")

	// Authentication
	apiv1.POST("/login", c.Auth.Login)
	apiv1.POST("/signup", c.Auth.SignUp)

	productRoutes := apiv1.Group("products")
	productRoutes.Use(c.Middleware.AuthMiddleware(), c.Middleware.StoreMiddleware())

	saleRoutes := apiv1.Group("sales")
	saleRoutes.Use(c.Middleware.AuthMiddleware(), c.Middleware.StoreMiddleware())

	userRoutes := apiv1.Group("user")
	userRoutes.Use(c.Middleware.AuthMiddleware())

	storeRoutes := apiv1.Group("store")
	storeRoutes.Use(c.Middleware.AuthMiddleware())

	productRouter(productRoutes, c)
	saleRouter(saleRoutes, c)
	storeRouter(storeRoutes, c)
	userRouter(userRoutes, c)

	return r
}

func productRouter(apiProduct *gin.RouterGroup, c controller.AppController) {

	apiProduct.POST("", c.Product.InsertProduct)
	apiProduct.GET("", c.Product.ListProduct)
	apiProduct.PUT("/:id", c.Product.UpdateProduct)
	apiProduct.DELETE("/:id", c.Product.DeleteProduct)
	apiProduct.GET("/search", c.Product.SearchProduct)

}

func storeRouter(apiStore *gin.RouterGroup, c controller.AppController) {

	apiStore.POST("", c.Employee.CreateStore)
	apiStore.POST("/:id", c.Employee.JoinStore)

	apiStore.Use(c.Middleware.StoreMiddleware(), c.Middleware.OwnerMiddleware())
	{
		apiStore.GET("", c.Employee.GetStoreInfo)
		apiStore.GET("/requests", c.Employee.GetJoinRequest)
		apiStore.PUT("/requests/:id", c.Employee.UpdateJoinRequest)
	}
}

func saleRouter(apiSale *gin.RouterGroup, c controller.AppController) {

	apiSale.POST("", c.Sale.CreateSale)
	apiSale.GET("/:id", c.Sale.GetSaleByID)
	apiSale.GET("", c.Sale.GetSales)
}

func userRouter(apiUser *gin.RouterGroup, c controller.AppController) {
	apiUser.Use(c.Middleware.AuthMiddleware())
	{
		apiUser.GET("", c.Employee.GetEmployeeInfo)
	}
}
