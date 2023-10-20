package registry

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/api"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
)

func (r *registry) NewProductController() api.ProductAPI {
	productRepo := respository.NewProductRepository(r.tx)
	productUseCases := usecases.NewProductUseCases(productRepo)
	p := api.NewProductAPI(
		productUseCases,
	)

	return p
}

func (r *registry) NewAuthController() api.AuthAPI {
	authRepo := respository.NewAuthRepositopry(r.tx)
	authUseCases := usecases.NewAuthUseCases(authRepo)
	a := api.NewAuthAPI(
		authUseCases,
	)

	return a
}

func (r *registry) NewSaleController() api.SaleAPI {
	saleRepo := respository.NewSaleRepository(r.tx)
	saleUseCases := usecases.NewSaleUseCases(saleRepo)
	a := api.NewSaleAPI(
		saleUseCases,
	)

	return a
}

func (r *registry) NewEmployeeController() api.EmployeeAPI {
	a := api.NewEmployeeAPI(
		respository.NewEmployeeRepository(r.tx),
	)

	return a
}
