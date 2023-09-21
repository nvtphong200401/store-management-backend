package respository

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	auth "github.com/nvtphong200401/store-management/pkg/handlers/models/auth"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SignUp(username string, password string) error
	Login(username string, password string) (map[string]string, error)
	CheckID(id uint, user *models.Employee) error
	VerifyToken(tokenString string) (*auth.Claims, error)
	RenewToken(refreshToken string) (string, error)
}

type authRepositoryImpl struct {
	tx *db.TxStore
}

func NewAuthRepositopry(tx *db.TxStore) AuthRepository {
	return &authRepositoryImpl{
		tx: tx,
	}
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

	claims := &auth.Claims{
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
func (r *authRepositoryImpl) RenewToken(refreshToken string) (string, error) {
	claims, err := r.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}
	if claims.Subject != helpers.REFRESH_TOKEN {
		return "", fmt.Errorf("invalid token type")
	}
	return generateToken(claims.UserID, helpers.ACCESS_TOKEN)
}

func (r *authRepositoryImpl) VerifyToken(tokenString string) (*auth.Claims, error) {
	jwtKey := os.Getenv(helpers.JWT_KEY)
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (r *authRepositoryImpl) Login(username string, password string) (map[string]string, error) {
	var user models.Employee
	// check username
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		result := db.Where("username = ?", username).First(&user)
		return result.Error
	})

	// match username
	if err != nil {
		return map[string]string{}, err
	}
	// check password
	err = comparePassword(string(user.Password), password)
	if err != nil {
		return map[string]string{}, err
	}
	// match password
	accessToken, err := generateToken(user.ID, helpers.ACCESS_TOKEN)
	if err != nil {
		return map[string]string{}, err
	}
	refreshToken, err := generateToken(user.ID, helpers.REFRESH_TOKEN)
	return map[string]string{helpers.ACCESS_TOKEN: accessToken, helpers.REFRESH_TOKEN: refreshToken}, err

}

func (r *authRepositoryImpl) SignUp(username string, password string) error {
	hashedPass, err := hashPassword(password)
	if err != nil {
		return err
	}
	var employee models.Employee = models.Employee{User: auth.User{Username: username, Password: hashedPass}, Position: models.Unknown}
	// create if not exist
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {

		db.AutoMigrate(&models.Employee{})

		return db.Create(&employee).Error
	})
}

func (r *authRepositoryImpl) CheckID(id uint, user *models.Employee) error {
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.First(&user, id).Error

	})
}
