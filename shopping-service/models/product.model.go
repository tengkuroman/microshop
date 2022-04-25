package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name, Description, ImageURL string
	Price                       int
	UserID, CategoryID          uint
}
