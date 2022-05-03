package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/tengkuroman/microshop/product-service/models"
	"github.com/tengkuroman/microshop/product-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 	Health check.
// @Description Connection health check.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1 [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "product",
	})
}

// @Summary 	Get all products.
// @Description Get all products available in marketplace.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/products [get]
func GetAllProducts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	if err := db.Find(&products).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var productsResponse []models.ProductResponse
	copier.Copy(&productsResponse, &products)

	response := utils.ResponseAPI("Get all products success!", http.StatusOK, "success", productsResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Get product by ID.
// @Description Get specific product by product_id.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/product/{product_id} [get]
// @Param 		product_id path int true "Param required."
func GetProductByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product

	if err := db.First(&product, c.Param("product_id")).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var productResponse models.ProductResponse
	copier.Copy(&productResponse, &product)

	response := utils.ResponseAPI("Get product success!", http.StatusOK, "success", productResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Get products from specific seller.
// @Description Get specific products by seller_id.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/products/seller/{user_id} [get]
// @Param 		user_id path int true "Param required."
func GetProductsBySellerID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	if err := db.Where("user_id = ?", c.Param("user_id")).Find(&products).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var productsResponse []models.ProductResponse
	copier.Copy(&productsResponse, &products)

	response := utils.ResponseAPI("Get products by seller ID success!", http.StatusOK, "success", productsResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Get products from specific category.
// @Description Get specific products by category_id.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/products/category/{category_id} [get]
// @Param 		category_id path int true "Param required."
func GetProductsByCategoryID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	if err := db.Where("category_id = ?", c.Param("category_id")).Find(&products).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var productsResponse []models.ProductResponse
	copier.Copy(&productsResponse, &products)

	response := utils.ResponseAPI("Get products by category ID success!", http.StatusOK, "success", productsResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Post product (role: seller)
// @Description Post product to marketplace. Switch your role if you are not seller.
// @Tags 		Product Service
// @Param 		body body models.ProductInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/product [post]
// @Security 	BearerToken
func PostProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	xUserID := c.Request.Header.Get("X-User-ID")

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", xUserID), 10, 32)
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		UserID:      uint(userID),
		CategoryID:  input.CategoryID,
	}

	if err := db.Create(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Product created successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Update product (role: seller)
// @Description Update posted product by product_id. Seller can only update their own products. Switch your role if you are not seller.
// @Tags 		Product Service
// @Param 		body body models.ProductInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/product/{product_id} [patch]
// @Param 		product_id path int true "Param required."
// @Security 	BearerToken
func UpdateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var product models.Product

	if err := db.Where("id = ?", c.Param("product_id")).First(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	productUserID := strconv.FormatUint(uint64(product.UserID), 10)
	userID := c.Request.Header.Get("X-User-ID")

	if productUserID != userID {
		response := utils.ResponseAPI("You can only update your own product!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var productInput models.ProductInput

	if err := c.ShouldBindJSON(&productInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := db.Model(&product).Updates(productInput).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Product data changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Delete product (role: seller)
// @Description Delete posted product by product_id. Seller can only delete their own products. Switch your role if you are not seller.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/product/{product_id} [delete]
// @Param 		product_id path int true "Param required."
// @Security 	BearerToken
func DeleteProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product

	if err := db.Where("id = ?", c.Param("product_id")).First(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	productUserID := strconv.FormatUint(uint64(product.UserID), 10)
	userID := c.Request.Header.Get("X-User-ID")

	if productUserID != userID {
		response := utils.ResponseAPI("You can only delete your own product!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Product deleted successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
