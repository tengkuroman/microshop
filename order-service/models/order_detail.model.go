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

type OrderDetailResponse struct {
	ID                uint                `json:"id"`
	Total             int                 `json:"total"`
	PaymentStatus     string              `json:"payment_status"`
	UserID            uint                `json:"user_id"`
	PaymentProviderID uint                `json:"payment_provider_id"`
	OrderItemResponse []OrderItemResponse `json:"order_item"`
}
