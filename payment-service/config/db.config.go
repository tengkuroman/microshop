package config

import (
	"fmt"
	"os"

	"github.com/tengkuroman/microshop/payment-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := os.Getenv("PAYMENT_DB_USERNAME")
	password := os.Getenv("PAYMENT_DB_PASSWORD")
	host := os.Getenv("PAYMENT_DB_HOST")
	port := os.Getenv("PAYMENT_DB_PORT")
	database := os.Getenv("PAYMENT_DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.PaymentProvider{})

	return db
}
