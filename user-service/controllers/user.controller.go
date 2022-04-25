package controllers

import (
	"net/http"
	"strconv"

	"github.com/tengkuroman/microshop/user-service/models"
	"github.com/tengkuroman/microshop/user-service/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var registerInput models.RegisterInput

	if err := c.ShouldBindJSON(&registerInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := models.User{
		FirstName:   registerInput.FirstName,
		LastName:    registerInput.LastName,
		Username:    registerInput.Username,
		Email:       registerInput.Email,
		Password:    registerInput.Password,
		Address:     registerInput.Address,
		PhoneNumber: registerInput.PhoneNumber,
		Role:        "user", // default role when registering
	}

	_, err := user.SaveUser(db)
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.ResponseAPI("Registration success!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var loginInput models.LoginInput

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := models.LoginCheck(loginInput.Username, loginInput.Password, db)
	if err != nil {
		response := utils.ResponseAPI("Username or password is incorrect!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := utils.ResponseAPI("Login success!", http.StatusOK, "success", map[string]string{"token": token})
	c.JSON(http.StatusOK, response)
}

func ChangePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var changePasswordInput models.ChangePasswordInput

	if err := c.ShouldBindJSON(&changePasswordInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userData, err := utils.ExtractPayload(c)
	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var user models.User

	if err := db.Where("id = ?", userData["user_id"]).First(&user).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := models.VerifyPassword(changePasswordInput.OldPassword, user.Password); err != nil {
		response := utils.ResponseAPI("Old password not match!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordInput.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response := utils.ResponseAPI("Password hashing error!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db.Model(&user).Update("password", newHashedPassword)

	response := utils.ResponseAPI("Password changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func ChangeUserDetail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var changeUserDetailInput models.ChangeUserDetailInput

	if err := c.ShouldBindJSON(&changeUserDetailInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var user models.User

	userData, err := utils.ExtractPayload(c)
	if err != nil {
		response := utils.ResponseAPI("Extract payload data failed!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if err := db.Where("id = ?", userData["user_id"]).First(&user).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Model(&user).Updates(changeUserDetailInput)

	response := utils.ResponseAPI("User details changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func SwitchUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	userData, err := utils.ExtractPayload(c)
	if err != nil {
		response := utils.ResponseAPI("Extract payload data failed!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if err := db.Where("id = ?", userData["user_id"]).First(&user).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newRole := c.Param("role")

	switch newRole {
	case "user":
	case "seller":
	case "admin":
	default:
		response := utils.ResponseAPI("Role invalid!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if newRole == user.Role {
		response := utils.ResponseAPI("Already in requested role!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Model(&user).Update("role", newRole)

	userInfo := map[string]string{
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"role":    user.Role,
	}

	token, err := utils.GenerateToken(userInfo)

	if err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseAPI("Role changed successfully!", http.StatusOK, "success", map[string]string{"token": token})
	c.JSON(http.StatusOK, response)
}

func ValidateUser(c *gin.Context) {
	err := utils.ValidateToken(c)
	if err != nil {
		response := utils.ResponseAPI("Token invalid!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userData, err := utils.ExtractPayload(c)
	if err != nil {
		response := utils.ResponseAPI("Extract payload data failed!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	authData := map[string]interface{}{
		"user_id": userData["user_id"],
		"role":    userData["role"],
	}

	response := utils.ResponseAPI("Token valid!", http.StatusOK, "success", authData)
	c.JSON(http.StatusOK, response)
}
