package controllers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/tengkuroman/microshop/shopping-service/models"
	"github.com/tengkuroman/microshop/shopping-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Connection to product service config
var (
	productHost    = os.Getenv("PRODUCT_HOST")
	productPort    = os.Getenv("PRODUCT_PORT")
	productBaseURL = fmt.Sprintf("%s:%s", productHost, productPort)
)

// Connection to product service config
var (
	orderHost    = os.Getenv("ORDER_HOST")
	orderPort    = os.Getenv("ORDER_PORT")
	orderBaseURL = fmt.Sprintf("%s:%s", orderHost, orderPort)
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "shopping",
	})
}

func AddProductToCart(c *gin.Context) {
	// Check active shopping session by user_id
	//      if exist then add to current session, update total in session
	//      if not exist then create session and add the product, update total in session

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	var item models.CartItem

	if err := c.ShouldBindJSON(&item); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := db.Where("user_id = ?", userData["user_id"]).Last(&session).Error; err != nil {
		UserID, err := strconv.ParseUint(fmt.Sprintf("%v", userData["user_id"]), 10, 32)
		if err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		session.UserID = uint(UserID)
		db.Create(&session)

		db.Where("user_id = ?", userData["user_id"]).Last(&session)

		item.ShoppingSessionID = session.ID
		db.Create(&item)

		productID := strconv.FormatUint(uint64(item.ProductID), 10)

		client := resty.New()
		res, err := client.R().SetResult(&models.ProductResponse{}).Get(productBaseURL + "/product/" + productID)

		if err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		product := res.Result().(*models.ProductResponse).Data

		itemTotalPrice := product.Price * item.Quantity
		db.Model(&session).Update("total", itemTotalPrice)

		response := utils.ResponseAPI("Product added to the cart!", http.StatusOK, "success", session)
		c.JSON(http.StatusOK, response)

		return
	}

	item.ShoppingSessionID = session.ID
	db.Create(&item)

	productID := strconv.FormatUint(uint64(item.ProductID), 10)

	client := resty.New()
	res, err := client.R().SetResult(&models.ProductResponse{}).Get(productBaseURL + "/product/" + productID)

	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	product := res.Result().(*models.ProductResponse).Data

	totalPrice := session.Total + item.Quantity*product.Price
	db.Model(&session).Update("total", totalPrice)

	response := utils.ResponseAPI("Product added to the cart!", http.StatusOK, "success", session)
	c.JSON(http.StatusOK, response)
}

func GetCartItems(c *gin.Context) {
	// Check active shopping session by user_id
	//      If exist then get items by shopping session
	//		If not exist then return "no items added to the cart"

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession

	if err := db.Where("user_id = ?", userData["user_id"]).Last(&session).Error; err != nil {
		response := utils.ResponseAPI("No items added to the cart!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var items []models.CartItem

	if err := db.Where("session_id = ?", session.ID).Find(&items).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.ResponseAPI("Get cart item success!", http.StatusOK, "success", items)
	c.JSON(http.StatusOK, response)
}

func UpdateCartItem(c *gin.Context) {
	// Check active shopping session by user_id -> get the shopping session id
	//     If exist then check if in shopping session there is a product_id == update item's product_id
	//			If exist then merge the data with the update
	//			If not exist then return "please use add product method"

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession

	if err := db.Where("user_id = ?", userData["user_id"]).Last(&session).Error; err != nil {
		response := utils.ResponseAPI("Please use add product method!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var updateItem models.CartItemUpdateInput

	if err := c.ShouldBindJSON(&updateItem); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var item models.CartItem

	if err := db.Where(&models.CartItem{ShoppingSessionID: session.ID, ProductID: updateItem.ProductID}).First(&item).Error; err != nil {
		response := utils.ResponseAPI("Please use add product method!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	itemOldQuantity := item.Quantity

	db.Model(&item).Update("quantity", updateItem.Quantity)

	productID := strconv.FormatUint(uint64(item.ProductID), 10)

	client := resty.New()
	res, err := client.R().SetResult(&models.ProductResponse{}).Get(productBaseURL + "/product/" + productID)

	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	product := res.Result().(*models.ProductResponse).Data

	// Set new total
	totalPrice := session.Total + product.Price*int(math.Abs(float64(item.Quantity-itemOldQuantity)))
	db.Model(&session).Update("total", totalPrice)

	response := utils.ResponseAPI("Cart item updated successfully!", http.StatusOK, "success", item)
	c.JSON(http.StatusOK, response)
}

func DropCart(c *gin.Context) {
	// Check active shopping session by user_id
	//		If exist then delete session and delete all cart items related to the session
	//		If not exist then return "no cart to be dropped"

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession

	if err := db.Where("user_id = ?", userData["user_id"]).Last(&session).Error; err != nil {
		response := utils.ResponseAPI("No cart to be dropped!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var item models.CartItem
	db.Where("session_id = ?", session.ID).Delete(&item)
	db.Delete(&session)

	response := utils.ResponseAPI("Cart dropped successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func Checkout(c *gin.Context) {
	// Check active shopping session by user_id
	//		If exist then:
	//			Create order detail and get key, get cart items and store to order item with order detail key
	//			Delete session and delete all cart items related to the session
	//		If not exist then return "no cart to be checked out"

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession

	if err := db.Where("user_id = ?", userData["user_id"]).Last(&session).Error; err != nil {
		response := utils.ResponseAPI("No cart to be checked out!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var cartItems []models.CartItem
	if err := db.Where("session_id = ?", session.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"session": session,
		"items":   cartItems,
	}

	client := resty.New()
	resp, err := client.R().SetBody(data).Post(orderBaseURL + "/order")

	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var modelCartItem models.CartItem
	db.Where("session_id = ?", session.ID).Delete(&modelCartItem)
	db.Delete(&session)

	response := utils.ResponseAPI("Order created successfully!", http.StatusOK, "success", resp.Body())
	c.JSON(http.StatusOK, response)
}
