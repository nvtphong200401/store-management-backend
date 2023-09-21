package db

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresDB(environmentPath string) (*gorm.DB, error) {
	err := godotenv.Load(environmentPath)

	if err != nil {
		return nil, err
	}
	host := os.Getenv(helpers.POSTGRES_HOST)
	port := os.Getenv(helpers.POSTGRES_PORT)
	user := os.Getenv(helpers.POSTGRES_USER)
	password := os.Getenv(helpers.POSTGRES_PASSWORD)
	dbname := os.Getenv(helpers.POSTGRES_DB)

	// Construct connection string
	dsn := fmt.Sprintf("host=%v port=%s user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbname)
	// dsn := "postgresql://phong:oV9rXZjpWTpD4yUHLb9Hyw@marsh-wren-4775.8nk.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func ConnectRedis(environmentPath string) *redis.Client {
	err := godotenv.Load(environmentPath)
	if err != nil {
		return nil
	}
	host := os.Getenv(helpers.REDIS_HOST)
	port := os.Getenv(helpers.REDIS_PORT)
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: "",
		DB:       0,
	})
}
