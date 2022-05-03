package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/tengkuroman/microshop/order-service/models"
	"github.com/tengkuroman/microshop/order-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Connection to payment service config
var (
	paymentHost    = os.Getenv("PAYMENT_HOST")
	paymentPort    = os.Getenv("PAYMENT_PORT")
	paymentBaseURL = fmt.Sprintf("%s:%s", paymentHost, paymentPort)
)

// @Summary 	Health check.
// @Description Connection health check.
// @Tags 		Order Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/order/v1 [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "order",
	})
}

// @Summary 	Get all user's order.
// @Description Get all user's order. Order retrieved only that made by logged user.
// @Tags 		Order Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/order/v1/orders [get]
// @Security 	BearerToken
func GetOrdersDetail(c *gin.Context) {
	// Get orders by user_id
	db := c.MustGet("db").(*gorm.DB)
	var orders []models.OrderDetail
	userID := c.Request.Header.Get("X-User-ID")

	db.Where("user_id = ?", userID).Find(&orders)

	var orderDetailResponse []models.OrderDetailResponse
	copier.Copy(&orderDetailResponse, &orders)

	response := utils.ResponseAPI("Get orders detail success!", http.StatusOK, "success", orderDetailResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Delete user's order.
// @Description Delete user's order. A user only can delete their own order.
// @Tags 		Order Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/order/v1/order/delete/{order_detail_id} [delete]
// @Param 		order_detail_id path int true "Param required."
// @Security 	BearerToken
func DeleteOrder(c *gin.Context) {
	// Check if an order exist based on param :order_detail_id
	// 		If order exist then check if order_detail.user.id == user_id
	//			OK: delete order_item where order_item.order_detail_id == order_detail.id, delete order_detail
	//			Not OK: Return message "You can only delete your order!"
	//		If order not exist then return "order detail not found"
	db := c.MustGet("db").(*gorm.DB)

	var order models.OrderDetail
	if err := db.Where("id = ?", c.Param("order_detail_id")).First(&order).Error; err != nil {
		response := utils.ResponseAPI("Order detail not found!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	orderUserID := strconv.FormatUint(uint64(order.UserID), 10)
	userID := c.Request.Header.Get("X-User-ID")

	if orderUserID == userID {
		var item models.OrderItem
		db.Where("order_detail_id = ?", order.ID).Delete(&item)
		db.Delete(&order)

		response := utils.ResponseAPI("Order deleted successfully!", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
	} else {
		response := utils.ResponseAPI("You can only delete your order!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
	}
}

// @Summary 	Select payment provider.
// @Description Select payment merchant after checkout (order created). A user can only pay their own order.
// @Tags 		Order Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/order/v1/order/payment/{order_detail_id} [patch]
// @Param 		order_detail_id path int true "Param required."
// @Security 	BearerToken
func SelectPaymentProvider(c *gin.Context) {
	// Check if an order exist based on param :order_detail_id
	// 		If order exist then check if order_detail.user.id == user_id
	//			OK: Update order_detail.payment_provider_id
	//			Not OK: Return message "You can only process payment of your order!"
	//		If order not exist then return "order detail not found"
	db := c.MustGet("db").(*gorm.DB)

	var order models.OrderDetail
	if err := db.Where("id = ?", c.Param("order_detail_id")).First(&order).Error; err != nil {
		response := utils.ResponseAPI("Order detail not found!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	orderUserID := strconv.FormatUint(uint64(order.UserID), 10)
	userID := c.Request.Header.Get("X-User-ID")

	if orderUserID == userID {
		db.Model(&order).Update("payment_provider_id", c.Param("payment_provider_id"))

		response := utils.ResponseAPI("Set payment provider success!", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
	} else {
		response := utils.ResponseAPI("You can only process payment of your order!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
	}
}

// @Summary 	Pay the order.
// @Description	Pay the selected order. A user can only pay their own order.
// @Tags 		Order Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/order/v1/order/payment/checkout/{order_detail_id} [patch]
// @Param 		order_detail_id path int true "Param required."
// @Security 	BearerToken
func PayOrder(c *gin.Context) {
	// Check if an order exist based on param :order_detail_id
	// 		If order exist then check if order_detail.user.id == user_id
	//			OK: If order paid?
	//					OK: Update order_detail.payment_status
	//					Not OK: Return message "Order already paid!"
	//			Not OK: Return message "You can only pay your order!"
	//		If order not exist then return "order detail not found"
	db := c.MustGet("db").(*gorm.DB)

	var order models.OrderDetail
	if err := db.Where("id = ?", c.Param("order_detail_id")).First(&order).Error; err != nil {
		response := utils.ResponseAPI("Order detail not found!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	orderUserID := strconv.FormatUint(uint64(order.UserID), 10)
	userID := c.Request.Header.Get("X-User-ID")

	if orderUserID == userID {
		if order.PaymentStatus == "paid" {
			response := utils.ResponseAPI("Order already paid!", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			data := map[string]interface{}{
				"total":               order.Total,
				"payment_provider_id": order.PaymentProviderID,
			}

			client := resty.New()
			resp, err := client.R().SetBody(data).Post(paymentBaseURL + "/payment/process")

			if err != nil {
				response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}

			fmt.Println(resp)

			db.Model(&order).Update("payment_status", "paid")

			response := utils.ResponseAPI("Order payment success!", http.StatusOK, "success", nil)
			c.JSON(http.StatusOK, response)
		}
	} else {
		response := utils.ResponseAPI("You can only process payment of your order!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
	}
}

// Invoked by shopping service
func CreateOrder(c *gin.Context) {
	// Bind session to order detail
	// Set payment status unpaid
	// Create to DB, get order detail ID
	// Create order item using order detail ID and items from REST
	var orderInput models.OrderInput

	if err := c.ShouldBindJSON(&orderInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	orderDetailInput := orderInput.Data.Session

	var orderDetail models.OrderDetail
	orderDetail.Total = orderDetailInput.Total
	orderDetail.PaymentStatus = "unpaid" //default when checkout
	orderDetail.UserID = orderDetailInput.UserID

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&orderDetail)

	orderItemsInput := orderInput.Data.Items

	var orderItems []models.OrderItem
	for i := range orderItemsInput {
		var orderItem models.OrderItem

		orderItem.Quantity = orderItemsInput[i].Quantity
		orderItem.ProductID = orderItemsInput[i].ProductID
		orderItem.OrderDetailID = orderDetail.ID

		orderItems = append(orderItems, orderItem)
	}

	db.Create(&orderItems)

	c.JSON(http.StatusOK, nil)
}
