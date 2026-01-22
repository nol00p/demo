package controllers

import (
	"demo/config"
	"demo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check for email
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password Invalid"})
		return
	}

	// Check for Password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password Invalid"})
		return
	}

	//if both email ans Password are correct
	claim := &CustomClaim{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte("a-string-secret-at-least-256-bits-long"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while generating the token "})
		return
	}

	c.JSON(http.StatusOK, tokenString)

}

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Allready in Use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User Creation Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Created"})
}
