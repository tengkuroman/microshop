package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name, Description, ImageURL string
	Price                       int
	UserID, CategoryID          uint
}

type ProductInput struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
	ImageURL    string `json:"image_url" binding:"required"`
	Price       int    `binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}
type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int    `json:"price"`
	UserID      uint   `json:"user_id"`
	CategoryID  uint   `json:"category_id"`
}
