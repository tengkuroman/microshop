package models

import (
	"html"
	"strconv"
	"strings"

	"github.com/tengkuroman/microshop/user-service/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `gorm:"not null;unique"`
	Email       string `gorm:"not null;unique"`
	Password    string
	Address     string
	PhoneNumber string `json:"phone_number"`
	Role        string
}

type RegisterInput struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Username    string `binding:"required"`
	Email       string `binding:"required"`
	Password    string `binding:"required"`
	Address     string
	PhoneNumber string `json:"phone_number"`
}

type LoginInput struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ChangeUserDetailInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string
	Address     string
	PhoneNumber string `json:"phone_number"`
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}

	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {
	var err error

	u := User{}

	err = db.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	userID := strconv.FormatUint(uint64(u.ID), 10)

	token, err := utils.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
