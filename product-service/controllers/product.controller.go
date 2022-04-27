package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tengkuroman/microshop/product-service/models"
	"github.com/tengkuroman/microshop/product-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "product",
	})
}

func GetAllProducts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	db.Find(&products)

	response := utils.ResponseAPI("Get all products success!", http.StatusOK, "success", products)

	c.JSON(http.StatusOK, response)
}

func GetProductByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product

	db.First(&product, c.Param("product_id"))

	response := utils.ResponseAPI("Get product success!", http.StatusOK, "success", product)

	c.JSON(http.StatusOK, response)
}

func GetProductsBySellerID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	db.Where("user_id = ?", c.Param("user_id")).Find(&products)

	response := utils.ResponseAPI("Get products by seller ID success!", http.StatusOK, "success", products)

	c.JSON(http.StatusOK, response)
}

func GetProductsByCategoryID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product

	db.Where("category_id = ?", c.Param("category_id")).Find(&products)

	response := utils.ResponseAPI("Get products by category ID success!", http.StatusOK, "success", products)

	c.JSON(http.StatusOK, response)
}

func PostProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", userData["user_id"]), 10, 32)
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

	db.Create(&product)

	response := utils.ResponseAPI("Product created successfully!", http.StatusOK, "success", product)
	c.JSON(http.StatusOK, response)
}

func UpdateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var product models.Product

	if err := db.Where("id = ?", c.Param("product_id")).First(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if product.UserID != userData["user_id"] {
		response := utils.ResponseAPI("You can only update your own product!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var input models.ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Model(&product).Updates(input)

	response := utils.ResponseAPI("Product data changed successfully!", http.StatusOK, "success", input)

	c.JSON(http.StatusOK, response)
}

func DeleteProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product

	if err := db.Where("id = ?", c.Param("product_id")).First(&product).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userData, err := utils.ExtractPayload(c.Query("token"))
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if product.UserID != userData["user_id"] {
		response := utils.ResponseAPI("You can only delete your own product!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db.Delete(&product)

	response := utils.ResponseAPI("Product deleted successfully!", http.StatusOK, "success", product)

	c.JSON(http.StatusOK, response)
}
