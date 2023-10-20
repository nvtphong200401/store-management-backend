package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
)

type AuthAPI interface {
	Login(c *gin.Context)
	SignUp(c *gin.Context)
	RenewToken(c *gin.Context)
}

type authAPIImpl struct {
	as usecases.AuthUseCases
}

func NewAuthAPI(ar usecases.AuthUseCases) AuthAPI {
	return &authAPIImpl{
		as: ar,
	}
}

func (api *authAPIImpl) RenewToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, "Refresh token is missing")
		return
	}

	statusCode, response := api.as.RenewToken(refreshToken)

	c.JSON(statusCode, response)
	// Return the new access token to the client
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, "Invalid refresh token")
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"access_token": token,
	// })
}

func (api *authAPIImpl) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	statusCode, response := api.as.Login(credentials.Username, credentials.Password)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid username or password"})
	// 	return
	// }

	c.JSON(statusCode, response)
}

func (api *authAPIImpl) SignUp(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	statusCode, response := api.as.SignUp(credentials.Username, credentials.Password)
	c.JSON(statusCode, response)
}
