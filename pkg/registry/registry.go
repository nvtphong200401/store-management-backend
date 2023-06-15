package registry

import (
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/controller"
	"gorm.io/gorm"
)

type Registry interface {
	NewAppController() controller.AppController
}
type registry struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewRegistry(db *gorm.DB, redisClient *redis.Client) Registry {
	return &registry{
		db:          db,
		redisClient: redisClient,
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
