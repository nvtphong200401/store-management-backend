package controller

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/api"
	"github.com/nvtphong200401/store-management/pkg/handlers/middleware"
)

type AppController struct {
	Product    api.ProductAPI
	Auth       api.AuthAPI
	Sale       api.SaleAPI
	Employee   api.EmployeeAPI
	Middleware middleware.Middleware
}
