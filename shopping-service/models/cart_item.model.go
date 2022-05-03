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

type CartItemUpdateInput struct {
	Quantity  int  `binding:"required"`
	ProductID uint `binding:"required"`
}

type CartItemResponse struct {
	ID                uint `json:"id"`
	Quantity          int  `json:"quantity"`
	ProductID         uint `json:"product_id"`
	ShoppingSessionID uint `json:"shopping_session_id"`
}
