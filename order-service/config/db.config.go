package config

import (
	"fmt"
	"os"

	"github.com/tengkuroman/microshop/order-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := os.Getenv("ORDER_DB_USERNAME")
	password := os.Getenv("ORDER_DB_PASSWORD")
	host := os.Getenv("ORDER_DB_HOST")
	port := os.Getenv("ORDER_DB_PORT")
	database := os.Getenv("ORDER_DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.OrderDetail{}, &models.OrderItem{})

	return db
}
