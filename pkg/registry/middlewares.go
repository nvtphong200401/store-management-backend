package registry

import (
	"github.com/nvtphong200401/store-management/pkg/handlers/middleware"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
)

func (r *registry) NewMiddleware() middleware.Middleware {
	return middleware.NewMiddleware(
		respository.NewAuthRepositopry(r.db),
	)
}
