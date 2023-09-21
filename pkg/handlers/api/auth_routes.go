package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
)

type AuthAPI interface {
	Login(c *gin.Context)
	SignUp(c *gin.Context)
	RenewToken(c *gin.Context)
}

type authAPIImpl struct {
	as respository.AuthRepository
}

func NewAuthAPI(ar respository.AuthRepository) AuthAPI {
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

	token, err := api.as.RenewToken(refreshToken)

	// Return the new access token to the client
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
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
	token, err := api.as.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
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
	if err := api.as.SignUp(credentials.Username, credentials.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
