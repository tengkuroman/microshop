package models

import "gorm.io/gorm"

type PaymentProvider struct {
	gorm.Model
	Name string
}

type PaymentProviderInput struct {
	Name string `binding:"required"`
}

type PaymentRequest struct {
	Total             int
	PaymentProviderID int
}

type PaymentProviderResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
