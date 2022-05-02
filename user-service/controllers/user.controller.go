package controllers

import (
	"net/http"

	"github.com/tengkuroman/microshop/user-service/models"
	"github.com/tengkuroman/microshop/user-service/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary 	Health check.
// @Description Connection health check.
// @Tags 		User Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/user/v1 [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Connection OK!",
		"service": "user",
	})
}

// @Summary 	Register a user.
// @Description Registering a user from public access.
// @Tags 		User Service
// @Param 		body body models.RegisterInput true "Body to register a user."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/user/v1/register [post]
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

// @Summary 	Login as as user, seller, or admin.
// @Description Logging in to get JWT token to access certain API by roles.
// @Tags 		User Service
// @Param 		body body models.LoginInput true "Body required to login."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/user/v1/login [post]
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

// @Summary 	Change user password.
// @Description Change user password for all roles.
// @Tags 		User Service
// @Param 		body body models.ChangePasswordInput true "Body required to change password."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/user/v1/change/password [patch]
// @Security 	BearerToken
func ChangePassword(c *gin.Context) {
	var changePasswordInput models.ChangePasswordInput

	if err := c.ShouldBindJSON(&changePasswordInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
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

// @Summary 	Change user details.
// @Description Change user detail: name, email, address, phone number.
// @Tags 		User Service
// @Param 		body body models.ChangeUserDetailInput true "Body required to user detail(s)."
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/user/v1/change [patch]
// @Security 	BearerToken
func ChangeUserDetail(c *gin.Context) {
	var changeUserDetailInput models.ChangeUserDetailInput

	if err := c.ShouldBindJSON(&changeUserDetailInput); err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		response := utils.ResponseAPI(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.Model(&user).Updates(changeUserDetailInput)

	response := utils.ResponseAPI("User details changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// @Summary 	Switch user role.
// @Description Change user role: user, seller, admin. Please re-login after switch.
// @Tags 		User Service
// @Produce 	json
// @Success 	200 {object} map[string]interface{}
// @Router 		/auth/user/v1/switch/{role} [patch]
// @Param 		role path string true "Available roles: user, seller, admin"
// @Security 	BearerToken
func SwitchUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	userID := c.Request.Header.Get("X-User-ID")

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
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

	response := utils.ResponseAPI("Role changed successfully!", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// Invoked by API gateway
func ValidateUser(c *gin.Context) {
	err := utils.ValidateToken(c)
	if err != nil {
		response := utils.ResponseAPI("Token invalid!", http.StatusUnauthorized, "unauthorized", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	claims, err := utils.ExtractPayload(c)
	if err != nil {
		response := utils.ResponseAPI("Extract payload data failed!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	err = db.Model(&user).Where("id = ?", claims["user_id"]).Take(&user).Error
	if err != nil {
		response := utils.ResponseAPI("Check user ID failed!", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"role":    user.Role,
	})
}
