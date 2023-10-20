package registry

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/middleware"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
)

func (r *registry) NewMiddleware() middleware.Middleware {
	authRepo := respository.NewAuthRepositopry(r.tx)
	authUseCases := usecases.NewAuthUseCases(authRepo)
	return middleware.NewMiddleware(
		authUseCases,
	)
}
