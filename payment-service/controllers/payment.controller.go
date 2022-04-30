package controllers

import (
	"net/http"

	"github.com/tengkuroman/microshop/payment-service/utils"

	"github.com/tengkuroman/microshop/payment-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "payment",
	})
}

func GetPaymentProviders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var providers []models.PaymentProvider

	db.Find(&providers)

	response := utils.ResponseAPI("Get payment providers success!", http.StatusOK, "success", providers)

	c.JSON(http.StatusOK, response)
}

func PostPaymentProvider(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can create payment provider!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var input models.PaymentProviderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	provider := models.PaymentProvider{
		Name: input.Name,
	}

	db.Create(&provider)

	response := utils.ResponseAPI("Payment provider created successfully!", http.StatusOK, "success", provider)

	c.JSON(http.StatusOK, response)
}

func UpdatePaymentProvider(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can update payment provider!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var input models.PaymentProviderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var provider models.PaymentProvider

	if err := db.Where("id = ?", c.Param("payment_provider_id")).First(&provider).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Model(&provider).Updates(input)

	response := utils.ResponseAPI("Payment provider changed successfully!", http.StatusOK, "success", provider)

	c.JSON(http.StatusOK, response)
}

func DeletePaymentProvider(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can delete payment provider!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var provider models.PaymentProvider

	if err := db.Where("id = ?", c.Param("payment_provider_id")).First(&provider).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Delete(&provider)

	response := utils.ResponseAPI("Payment provider deleted successfully!", http.StatusOK, "success", provider)

	c.JSON(http.StatusOK, response)
}

func ProcessPayment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var paymentRequest models.PaymentRequest

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var provider models.PaymentProvider

	if err := db.Where("id = ?", paymentRequest.PaymentProviderID).First(&provider).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.ResponseAPI("Payment processed successfully!", http.StatusOK, "success", paymentRequest)

	c.JSON(http.StatusOK, response)
}
