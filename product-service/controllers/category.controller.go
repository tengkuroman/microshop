package controllers

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/tengkuroman/microshop/product-service/models"
	"github.com/tengkuroman/microshop/product-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 	Get all product categories.
// @Description Get all product categories, including unsigned to product categories.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/categories [get]
func GetAllCategories(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var categories []models.Category

	if err := db.Find(&categories).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var categoriesResponse []models.CategoryResponse
	copier.Copy(&categoriesResponse, &categories)

	response := utils.ResponseAPI("Get all categories success!", http.StatusOK, "success", categoriesResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Get product category by ID.
// @Description Get product category by category_id.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/product/v1/category/{category_id} [get]
// @Param 		category_id path string true "Required param."
func GetCategoryByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.First(&category, c.Param("category_id")).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var categoryResponse models.CategoryResponse
	copier.Copy(&categoryResponse, &category)

	response := utils.ResponseAPI("Get category success!", http.StatusOK, "success", categoryResponse)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Post category (role: admin)
// @Description Post product category. Only admin can post category. Switch your role if you are not admin.
// @Tags 		Product Service
// @Param 		body body models.CategoryInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/category [post]
// @Security 	BearerToken
func PostCategory(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can post category!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var input models.CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := db.Create(&category).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Category created successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Update product category (role: admin)
// @Description Update product category by category_id. Only admin can update category. Switch your role if you are not admin.
// @Tags 		Product Service
// @Param 		body body models.CategoryInput true "Body required."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/category/{category_id} [patch]
// @Param 		category_id path int true "Param required."
// @Security 	BearerToken
func UpdateCategory(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can update category!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var categoryInput models.CategoryInput

	if err := c.ShouldBindJSON(&categoryInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var category models.Category

	if err := db.Where("id = ?", c.Param("category_id")).First(&category).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := db.Model(&category).Updates(categoryInput).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var categoryResponse models.CategoryResponse
	copier.Copy(&categoryResponse, &category)

	response := utils.ResponseAPI("Category data changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Delete product category (role: admin)
// @Description Delete product category. Only admin can delete category. Switch your role if you are not admin.
// @Tags 		Product Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/product/v1/category/{category_id} [delete]
// @Param 		category_id path int true "Param required."
// @Security 	BearerToken
func DeleteCategory(c *gin.Context) {
	userRole := c.Request.Header.Get("X-User-Role")

	if userRole != "admin" {
		response := utils.ResponseAPI("Only admins can delete category!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.Where("id = ?", c.Param("category_id")).First(&category).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Category deleted successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
