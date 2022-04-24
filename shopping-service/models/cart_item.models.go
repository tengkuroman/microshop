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
