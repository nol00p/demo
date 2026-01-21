package controllers

import (
	"demo/config"
	"demo/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	if err := config.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get Project data"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := config.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project can't be found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func PostProject(c *gin.Context) {
	var project models.Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating Project"})
	}

	c.JSON(http.StatusOK, project)
}
