package usecases

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCases struct {
	authRepo respository.AuthRepository
}

type Claims struct {
	UserID uint // Add additional claims as needed
	jwt.StandardClaims
}

func NewAuthUseCases(authRepo respository.AuthRepository) AuthUseCases {
	return AuthUseCases{
		authRepo: authRepo,
	}
}

func (a *AuthUseCases) Login(username string, password string) (int, gin.H) {
	user, err := a.authRepo.GetUserByUsername(username)
	if err != nil {
		return responseWithError(err)
	}
	err = comparePassword(string(user.Password), password)
	if err != nil {
		return responseWithError(err)
	}

	accessToken, err := generateToken(user.ID, helpers.ACCESS_TOKEN)
	if err != nil {
		return responseWithError(err)
	}
	refreshToken, err := generateToken(user.ID, helpers.REFRESH_TOKEN)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(map[string]string{
		helpers.ACCESS_TOKEN:  accessToken,
		helpers.REFRESH_TOKEN: refreshToken,
	})
}

func (a *AuthUseCases) RenewToken(refreshToken string) (int, gin.H) {
	claims, err := a.VerifyToken(refreshToken)
	if err != nil {
		return responseWithError(err)
	}
	if claims.Subject != helpers.REFRESH_TOKEN {
		return responseWithError(fmt.Errorf("invalid token type"))
	}
	newAccessToken, err := generateToken(claims.UserID, helpers.ACCESS_TOKEN)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(newAccessToken)
}

func (a *AuthUseCases) GetUserInfo(tokenString string) (*models.Employee, error) {
	claims, err := a.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return a.authRepo.GetUserByID(claims.UserID)
}

func (a *AuthUseCases) VerifyToken(tokenString string) (*Claims, error) {
	jwtKey := os.Getenv(helpers.JWT_KEY)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// Sign up use
func (a *AuthUseCases) SignUp(username string, password string) (int, gin.H) {
	hashedPass, err := hashPassword(password)
	if err != nil {
		return responseWithError(err)
	}
	err = a.authRepo.SignUp(username, hashedPass)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithResult("success")
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateToken(userID uint, tokenType string) (string, error) {
	var expiredTime int64
	jwtKey := os.Getenv(helpers.JWT_KEY)
	switch tokenType {
	case helpers.ACCESS_TOKEN:
		expiredTime = time.Now().Add(time.Hour * 24).Unix()
	case helpers.REFRESH_TOKEN:
		expiredTime = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime,
			IssuedAt:  time.Now().Unix(),
			Subject:   tokenType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
