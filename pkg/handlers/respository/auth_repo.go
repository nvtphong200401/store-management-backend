package respository

import (
	"github.com/go-redis/redis"
	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	auth "github.com/nvtphong200401/store-management/pkg/handlers/models/auth"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SignUp(username string, password []byte) error
	GetUserByUsername(username string) (*models.Employee, error)
	GetUserByID(id uint) (*models.Employee, error)
}

type authRepositoryImpl struct {
	tx *db.TxStore
}

func NewAuthRepositopry(tx *db.TxStore) AuthRepository {
	return &authRepositoryImpl{
		tx: tx,
	}
}

func (r *authRepositoryImpl) GetUserByUsername(username string) (*models.Employee, error) {
	var user models.Employee
	// check username
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		result := db.Where("username = ?", username).First(&user)
		return result.Error
	})

	// match username
	if err != nil {
		return nil, err
	}
	return &user, nil
	// // check password
	// err = comparePassword(string(user.Password), password)
	// if err != nil {
	// 	return "", "", err
	// }
	// // match password
	// accessToken, err := generateToken(user.ID, helpers.ACCESS_TOKEN)
	// if err != nil {
	// 	return "", "", err
	// }
	// refreshToken, err := generateToken(user.ID, helpers.REFRESH_TOKEN)
	// return accessToken, refreshToken, err

}

func (r *authRepositoryImpl) SignUp(username string, password []byte) error {

	var employee models.Employee = models.Employee{User: auth.User{Username: username, Password: password}, Position: models.Unknown}
	// create if not exist
	return r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {

		db.AutoMigrate(&models.Employee{})

		return db.Create(&employee).Error
	})
}

func (r *authRepositoryImpl) GetUserByID(id uint) (*models.Employee, error) {
	var user models.Employee
	err := r.tx.ExecuteTX(func(db *gorm.DB, rd *redis.Client) error {
		return db.First(&user, id).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}
