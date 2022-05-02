package controllers

import (
	"net/http"

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
// @Description Get payment providers. Authenticated user can access this.
// @Tags 		Payment Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/payment/v1/payment [get]
func GetPaymentProviders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var providers []models.PaymentProvider

	db.Find(&providers)

	response := utils.ResponseAPI("Get payment providers success!", http.StatusOK, "success", providers)
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

	db.Delete(&provider)

	response := utils.ResponseAPI("Payment provider deleted successfully!", http.StatusOK, "success", provider)
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
