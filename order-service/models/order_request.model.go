package models

// Model for service invocation from shopping service
type CartItemInput struct {
	Quantity  int  `binding:"required" json:"quantity"`
	ProductID uint `json:"product_id" binding:"required"`
}

type ShoppingSessionInput struct {
	Total  int  `binding:"required" json:"total"`
	UserID uint `json:"user_id" binding:"required"`
}

type OrderInput struct {
	Session ShoppingSessionInput `binding:"required" json:"session"`
	Items   []CartItemInput      `binding:"required" json:"items"`
}
