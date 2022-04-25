package models

type CartItemInput struct {
	Quantity  int  `binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}

type ShoppingSessionInput struct {
	Total  int  `binding:"required"`
	UserID uint `json:"user_id" binding:"required"`
}

type OrderInput struct {
	Data struct {
		Session ShoppingSessionInput `binding:"required"`
		Items   []CartItemInput      `binding:"required"`
	} `binding:"required"`
}
