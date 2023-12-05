package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
)

type AuthAPI interface {
	Login(c *gin.Context)
	SignUp(c *gin.Context)
	RenewToken(c *gin.Context)
	VerifyCode(c *gin.Context)
	RequestVerificationCode(c *gin.Context)
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

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @BasePath /api/v1

// Login godoc
// @Summary Login with email and password
// @Schemes
// @Description Login with email and password
// @Tags Login
// @Accept json
// @Produce json
// @Param			auth	body		credentials	true	"Auth"
// @Success 200 {object} credentials
// @Router /login [post]
func (api *authAPIImpl) Login(c *gin.Context) {
	var credentials credentials
	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	statusCode, response := api.as.Login(credentials.Email, credentials.Password)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
	// 	return
	// }

	c.JSON(statusCode, response)
}

func (api *authAPIImpl) SignUp(c *gin.Context) {
	var credentials credentials
	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}
	statusCode, response := api.as.SignUp(credentials.Email, credentials.Password)
	c.JSON(statusCode, response)
}

func (api *authAPIImpl) VerifyCode(c *gin.Context) {
	code := c.GetString("code")
	user, exist := c.Get("user")
	if exist {
		api.as.VerifyCode(user.(models.User), code)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("something went wrong").Error())
	}
}

func (api *authAPIImpl) RequestVerificationCode(c *gin.Context) {
	user, exist := c.Get("user")

	if exist {
		api.as.RequestVerificationCode(user.(models.User))
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("something went wrong").Error())
	}
}
