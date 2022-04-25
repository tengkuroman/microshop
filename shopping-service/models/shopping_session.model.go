package models

import "gorm.io/gorm"

type ShoppingSession struct {
	gorm.Model
	Total    int
	UserID   uint
	CartItem []CartItem
}
