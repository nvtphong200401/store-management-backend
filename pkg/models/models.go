package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Construct connection string
	dsn := fmt.Sprintf("host=%v port=%s user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

}

func CloseDB() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
