package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	Quantity          int
	ProductID         uint
	ShoppingSessionID uint
}

type CartItemInput struct {
	Quantity  int  `binding:"required" json:"quantity"`
	ProductID uint `binding:"required" json:"product_id"`
}

type CartItemResponse struct {
	ID                uint `json:"id"`
	Quantity          int  `json:"quantity"`
	ProductID         uint `json:"product_id"`
	ShoppingSessionID uint `json:"shopping_session_id"`
}
