package db

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client

func SetUp() (*gorm.DB, *redis.Client) {
	if db != nil && redisClient != nil {
		return db, redisClient
	}
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
	// dsn := "postgresql://phong:oV9rXZjpWTpD4yUHLb9Hyw@marsh-wren-4775.8nk.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return db, rd
}
