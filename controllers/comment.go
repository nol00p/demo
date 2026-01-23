package controllers

import (
	"demo/config"
	"demo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	userIDInt, ok := userID.(int)
	comment.UserID = uint(userIDInt)

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment can't be saved"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
