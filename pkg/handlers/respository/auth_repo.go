package respository

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/handlers/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SignUp(username string, password string) error
	Login(username string, password string) (string, error)
	CheckID(id uint, user *models.Employee) error
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
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("my-secret-key"))
}

func (r *authRepositoryImpl) Login(username string, password string) (string, error) {
	var user models.Employee
	// check username
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		result := db.Where("username = ?", username).First(&user)
		return result.Error
	})

	// match username
	if err != nil {
		return "", err
	}
	// check password
	err = comparePassword(string(user.Password), password)
	if err != nil {
		return "", err
	}
	// match password
	token, err := generateToken(user.ID)
	return token, nil
}

func (r *authRepositoryImpl) SignUp(username string, password string) error {
	hashedPass, err := hashPassword(password)
	if err != nil {
		return err
	}
	var employee models.Employee = models.Employee{User: models.User{Username: username, Password: hashedPass}, Position: models.Unknown}
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
