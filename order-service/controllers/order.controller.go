package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
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

func GetOrdersDetail(c *gin.Context) {
	// Get orders by user_id
	db := c.MustGet("db").(*gorm.DB)

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var orders []models.OrderDetail
	db.Where("user_id", userData["user_id"]).Find(&orders)

	response := utils.ResponseAPI("Get orders detail success!", http.StatusOK, "success", orders)
	c.JSON(http.StatusOK, response)
}

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

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if order.UserID == userData["user_id"] {
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

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if order.UserID == userData["user_id"] {
		db.Model(&order).Update("payment_provider_id", c.Param("payment_provider_id"))

		response := utils.ResponseAPI("Set payment provider success!", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
	} else {
		response := utils.ResponseAPI("You can only process payment of your order!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
	}
}

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

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if order.UserID == userData["user_id"] {
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

func CreateOrder(c *gin.Context) {
	// Bind session to order detail
	// Set payment status unpaid
	// Create to DB, get order detail ID
	// Create order item using order detail ID and items from REST
	db := c.MustGet("db").(*gorm.DB)

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

	response := utils.ResponseAPI("Create order success!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
