package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	Quantity      int
	ProductID     uint
	OrderDetailID uint
}

type OrderItemResponse struct {
	ID        uint `json:"id"`
	Quantity  uint `json:"quantity"`
	ProductID uint `json:"product_id"`
}
