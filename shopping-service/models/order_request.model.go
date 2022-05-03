package models

// Model for service invocation to order service
type CartItemOrder struct {
	Quantity  int  `binding:"required" json:"quantity"`
	ProductID uint `json:"product_id" binding:"required"`
}

type ShoppingSessionOrder struct {
	Total  int  `binding:"required" json:"total"`
	UserID uint `json:"user_id" binding:"required"`
}

type Order struct {
	Session ShoppingSessionOrder `binding:"required" json:"session"`
	Items   []CartItemInput      `binding:"required" json:"items"`
}
