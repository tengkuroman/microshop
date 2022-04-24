package controllers

import (
	"net/http"

	"github.com/tengkuroman/microshop/product-service/models"
	"github.com/tengkuroman/microshop/product-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllCategories(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var categories []models.Category

	db.Find(&categories)

	response := utils.ResponseAPI("Get all categories success!", http.StatusOK, "success", categories)

	c.JSON(http.StatusOK, response)
}

func GetCategoryByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	db.First(&category, c.Param("category_id"))

	response := utils.ResponseAPI("Get category success!", http.StatusOK, "success", category)

	c.JSON(http.StatusOK, response)
}

func PostCategory(c *gin.Context) {
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

	db.Create(&category)

	response := utils.ResponseAPI("Category created successfully!", http.StatusOK, "success", category)

	c.JSON(http.StatusOK, response)
}

func UpdateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input models.CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
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

	db.Model(&category).Updates(input)

	response := utils.ResponseAPI("Category data changed successfully!", http.StatusOK, "success", input)

	c.JSON(http.StatusOK, response)
}

func DeleteCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category

	if err := db.Where("id = ?", c.Param("category_id")).First(&category).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Delete(&category)

	response := utils.ResponseAPI("Category deleted successfully!", http.StatusOK, "success", category)

	c.JSON(http.StatusOK, response)
}
