package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password []byte
}

type AuthService struct{}

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

func (as *AuthService) Login(username string, password string) (string, error) {
	var user User
	// check username
	result := db.Where("username = ?", username).First(&user)
	// match username
	if result.Error != nil {
		return "", result.Error
	}
	// check password
	err := comparePassword(string(user.Password), password)
	if err != nil {
		return "", err
	}
	// match password
	token, err := generateToken(user.ID)
	return token, nil
}

func (as *AuthService) SignUp(username string, password string) error {
	hashedPass, err := hashPassword(password)
	if err != nil {
		return err
	}
	var user User = User{Username: username, Password: hashedPass}
	// create if not exist
	db.AutoMigrate(&User{})

	result := db.Create(&user)
	return result.Error
}

func (as *AuthService) CheckID(id uint) error {
	var user User
	result := db.First(&user, id)
	return result.Error
}
