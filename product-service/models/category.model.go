package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string
	Description string
	Product     []Product
}

type CategoryInput struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
}
