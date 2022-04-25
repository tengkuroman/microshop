package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	Quantity      int
	ProductID     uint
	OrderDetailID uint
}
