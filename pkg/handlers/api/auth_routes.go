package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/models"
	"golang.org/x/time/rate"
)

var (
	as         *models.AuthService
	limiter    *rate.Limiter
	blockedIPs = make(map[string]time.Time)
)

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	token, err := as.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func SignUp(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	if err := as.SignUp(credentials.Username, credentials.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cannot create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
