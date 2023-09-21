package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	StoreMiddleware() gin.HandlerFunc
	OwnerMiddleware() gin.HandlerFunc
}

type middlewareImpl struct {
	auth respository.AuthRepository
}

func NewMiddleware(a respository.AuthRepository) Middleware {
	return &middlewareImpl{
		auth: a,
	}
}

func (m *middlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, err := m.auth.VerifyToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		if err := claims.Valid(); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var user models.Employee

		m.auth.CheckID(claims.UserID, &user)
		c.Set("user", user)

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
