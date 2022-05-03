package controllers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
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

// @Summary 	Health check.
// @Description Connection health check.
// @Tags 		Shopping Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/shopping/v1 [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "shopping",
	})
}

// @Summary 	Add a product to cart.
// @Description Add a product to cart.
// @Tags 		Shopping Service
// @Param 		body body models.CartItemInput true "Body to add product to the cart."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/shopping/v1/cart [post]
// @Security 	BearerToken
func AddProductToCart(c *gin.Context) {
	// Check active shopping session by user_id
	//      if exist then add to current session, update total in session
	//      if not exist then create session and add the product, update total in session
	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	var item models.CartItem
	var itemInput models.CartItemInput

	if err := c.ShouldBindJSON(&itemInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	item.Quantity = itemInput.Quantity
	item.ProductID = itemInput.ProductID

	xUserID := c.Request.Header.Get("X-User-ID")
	if err := db.Where("user_id = ?", xUserID).Last(&session).Error; err != nil {
		userID, err := strconv.ParseUint(fmt.Sprintf("%v", xUserID), 10, 32)
		if err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		session.UserID = uint(userID)
		if err := db.Create(&session).Error; err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		item.ShoppingSessionID = session.ID
		if err := db.Create(&item).Error; err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		productID := strconv.FormatUint(uint64(item.ProductID), 10)
		client := resty.New()
		res, err := client.R().SetResult(&models.ProductResponse{}).Get("http://" + productBaseURL + "/product/" + productID)

		if err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		product := res.Result().(*models.ProductResponse).Data

		itemTotalPrice := product.Price * item.Quantity
		if err := db.Model(&session).Update("total", itemTotalPrice).Error; err != nil {
			response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := utils.ResponseAPI("Product added to the cart!", http.StatusOK, "success", nil)
		c.JSON(http.StatusOK, response)
		return
	}

	item.ShoppingSessionID = session.ID
	if err := db.Create(&item).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	productID := strconv.FormatUint(uint64(item.ProductID), 10)

	client := resty.New()
	res, err := client.R().SetResult(&models.ProductResponse{}).Get("http://" + productBaseURL + "/product/" + productID)

	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	product := res.Result().(*models.ProductResponse).Data

	totalPrice := session.Total + item.Quantity*product.Price
	if err := db.Model(&session).Update("total", totalPrice).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Product added to the cart!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Get all products from cart.
// @Description Get all products from cart. Data retrieved based on logged in user.
// @Tags 		Shopping Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/shopping/v1/cart [get]
// @Security 	BearerToken
func GetCartItems(c *gin.Context) {
	// Check active shopping session by user_id
	//      If exist then get items by shopping session
	//		If not exist then return "no items added to the cart"
	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("user_id = ?", userID).Last(&session).Error; err != nil {
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

	var itemsResponse []models.CartItemResponse
	copier.Copy(&itemsResponse, &items)

	response := utils.ResponseAPI("Get cart item success!", http.StatusOK, "success", itemsResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Update a product quantity in cart.
// @Description Update a product quantity in cart.
// @Tags 		Shopping Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/shopping/v1/cart [patch]
// @Security 	BearerToken
func UpdateCartItem(c *gin.Context) {
	// Check active shopping session by user_id -> get the shopping session id
	//     If exist then check if in shopping session there is a product_id == update item's product_id
	//			If exist then merge the data with the update
	//			If not exist then return "please use add product method"
	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("user_id = ?", userID).Last(&session).Error; err != nil {
		response := utils.ResponseAPI("Please use add product method!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var updateItem models.CartItemInput

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

	response := utils.ResponseAPI("Cart item updated successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Drop shopping cart.
// @Description Delete shopping session and all items in cart for current logged in user.
// @Tags 		Shopping Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/shopping/v1/cart [delete]
// @Security 	BearerToken
func DropCart(c *gin.Context) {
	// Check active shopping session by user_id
	//		If exist then delete session and delete all cart items related to the session
	//		If not exist then return "no cart to be dropped"
	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("user_id = ?", userID).Last(&session).Error; err != nil {
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

// @Summary 	Checkout shopping cart.
// @Description Bring all the items in cart to order.
// @Tags 		Shopping Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/shopping/v1/cart/checkout [get]
// @Security 	BearerToken
func Checkout(c *gin.Context) {
	// Check active shopping session by user_id
	//		If exist then:
	//			Create order detail and get key, get cart items and store to order item with order detail key
	//			Delete session and delete all cart items related to the session
	//		If not exist then return "no cart to be checked out"
	db := c.MustGet("db").(*gorm.DB)
	var session models.ShoppingSession
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("user_id = ?", userID).Last(&session).Error; err != nil {
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

	fmt.Println(resp)

	var modelCartItem models.CartItem
	db.Where("session_id = ?", session.ID).Delete(&modelCartItem)
	db.Delete(&session)

	response := utils.ResponseAPI("Order created successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
