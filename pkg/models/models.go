package models

import (
	"log"

	"github.com/nvtphong200401/store-management/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// type Model struct {
// 	ID        int       `gorm:"primary_key" json:"id"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"deleted_at"`
// }

func SetUp() {
	var err error
	db, err = gorm.Open(postgres.Open(config.DbInfo), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

}

func CloseDB() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
