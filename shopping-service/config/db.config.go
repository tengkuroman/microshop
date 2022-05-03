package config

import (
	"fmt"
	"os"

	"github.com/tengkuroman/microshop/shopping-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := os.Getenv("SHOPPING_DB_USERNAME")
	password := os.Getenv("SHOPPING_DB_PASSWORD")
	host := os.Getenv("SHOPPING_DB_HOST")
	port := os.Getenv("SHOPPING_DB_PORT")
	database := os.Getenv("SHOPPING_DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(
		&models.ShoppingSession{},
		&models.CartItem{},
	)

	return db
}
