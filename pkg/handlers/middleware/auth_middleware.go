package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

func verifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve the JWT secret from a configuration file or environment variable
		return []byte("my-secret-key"), nil
	})

}

func (m *middlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		token, err := verifyToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			userID := claims["user_id"].(float64)
			var user models.Employee
			err = m.auth.CheckID(uint(userID), &user)
			if err := claims.Valid(); err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Set the user information in the Gin context for downstream handlers to access
			c.Set("user", user)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

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
