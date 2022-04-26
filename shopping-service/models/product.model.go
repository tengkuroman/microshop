package models

type ProductResponse struct {
	Data struct {
		Name        string `binding:"required"`
		Description string `binding:"required"`
		ImageURL    string `json:"image_url" binding:"required"`
		Price       int    `binding:"required"`
		UserID      uint   `json:"user_id" binding:"required"`
		CategoryID  uint   `json:"category_url" binding:"required"`
	} `binding:"required"`
}
