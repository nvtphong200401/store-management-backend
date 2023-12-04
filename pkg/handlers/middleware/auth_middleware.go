package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	StoreMiddleware() gin.HandlerFunc
	OwnerMiddleware() gin.HandlerFunc
}

type middlewareImpl struct {
	authUseCase usecases.AuthUseCases
}

func NewMiddleware(a usecases.AuthUseCases) Middleware {
	return &middlewareImpl{
		authUseCase: a,
	}
}

func (m *middlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		user, err := m.authUseCase.GetUserInfo(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		// TODO: handle unverified user
		c.Set("user", *user)

		c.Next()
	}
}

func (m *middlewareImpl) StoreMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		employee, err := helpers.GetEmployee(c)
		if err == nil {
			if !employee.AlreadyInStore() {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "You have not created or joined a store"})
				return
			} else {
				c.Next()
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
	}
}

func (m *middlewareImpl) OwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		employee, err := helpers.GetEmployee(c)
		if err != nil {
			return
		}
		if employee.IsEmployeeOwner() {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		}
	}
}
