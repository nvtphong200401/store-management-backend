package registry

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/controller"
	"github.com/nvtphong200401/store-management/pkg/handlers/db"
)

type Registry interface {
	NewAppController() controller.AppController
}
type registry struct {
	tx *db.TxStore
}

func NewRegistry(tx db.TxStore) Registry {
	return &registry{
		tx: &tx,
	}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		Product:    r.NewProductController(),
		Auth:       r.NewAuthController(),
		Sale:       r.NewSaleController(),
		Employee:   r.NewEmployeeController(),
		Middleware: r.NewMiddleware(),
	}
}
