package registry

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/api"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
)

func (r *registry) NewProductController() api.ProductAPI {
	p := api.NewProductAPI(
		respository.NewProductRepository(r.db, r.redisClient),
	)

	return p
}

func (r *registry) NewAuthController() api.AuthAPI {
	a := api.NewAuthAPI(
		respository.NewAuthRepositopry(r.db),
	)

	return a
}

func (r *registry) NewSaleController() api.SaleAPI {
	a := api.NewSaleAPI(
		respository.NewSaleRepository(r.db),
	)

	return a
}

func (r *registry) NewEmployeeController() api.EmployeeAPI {
	a := api.NewEmployeeAPI(
		respository.NewEmployeeRepository(r.db, r.redisClient),
	)

	return a
}
