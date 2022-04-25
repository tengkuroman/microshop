package models

import "gorm.io/gorm"

type OrderDetail struct {
	gorm.Model
	Total             int
	PaymentStatus     string
	UserID            uint
	PaymentProviderID uint
	OrderItem         []OrderItem
}
