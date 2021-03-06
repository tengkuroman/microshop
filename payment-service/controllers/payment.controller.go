package controllers

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/tengkuroman/microshop/payment-service/utils"

	"github.com/tengkuroman/microshop/payment-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 	Health check.
// @Description Connection health check.
// @Tags 		Payment Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/payment/v1 [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "payment",
	})
}

// @Summary 	Get payment providers.
// @Description Get payment providers.
// @Tags 		Payment Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/payment/v1/payment [get]
func GetPaymentProviders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var providers []models.PaymentProvider

	if err := db.Find(&providers).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var providersResponse []models.PaymentProviderResponse
	copier.Copy(&providersResponse, &providers)

	response := utils.ResponseAPI("Get payment providers success!", http.StatusOK, "success", providersResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Post payment provider (role: admin)
// @Description Post payment provider. Only admin can post it. Switch your role if you are not admin.
// @Tags 		Payment Service
// @Param 		body body models.PaymentProviderInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/payment/v1/payment [post]
// @Security 	BearerToken
func PostPaymentProvider(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can create payment provider!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var paymentProviderInput models.PaymentProviderInput

	if err := c.ShouldBindJSON(&paymentProviderInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	provider := models.PaymentProvider{
		Name: paymentProviderInput.Name,
	}

	if err := db.Create(&provider).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Payment provider created successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Update payment provider (role: admin)
// @Description Update payment provider. Only admin update it. Switch your role if you are not admin.
// @Tags 		Payment Service
// @Param 		body body models.PaymentProviderInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/payment/v1/payment/{payment_provider_id} [patch]
// @Param 		payment_provider_id path int true "Param required."
// @Security 	BearerToken
func UpdatePaymentProvider(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can update payment provider!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var paymentProviderInput models.PaymentProviderInput

	if err := c.ShouldBindJSON(&paymentProviderInput); err != nil {
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

	if err := db.Model(&provider).Updates(paymentProviderInput).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Payment provider changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Delete payment provider (role: admin)
// @Description Delete payment provider. Only admin can delete it. Switch your role if you are not admin.
// @Tags 		Payment Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/payment/v1/payment/{payment_provider_id} [delete]
// @Param 		payment_provider_id path int true "Param required."
// @Security 	BearerToken
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

	if err := db.Delete(&provider).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Payment provider deleted successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// Invoked by order service
func ProcessPayment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var paymentRequest models.PaymentRequest

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var provider models.PaymentProvider

	if err := db.Where("id = ?", paymentRequest.PaymentProviderID).First(&provider).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payment provider not found!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Payment processed successfully!",
	})
}
