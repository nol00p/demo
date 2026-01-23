package controllers

import (
	"demo/config"
	"demo/models"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GetProjects godoc
// @Description Get all Projects
// @Tags Projects
// @Produce json
// @Success 200  {arrray} models.GetProjects
// @Security BearerAuth
// @Router /projects [get]
func GetProjects(c *gin.Context) {
	var projects []models.Project

	if err := config.DB.Preload("Comments").Preload("Likes").Find(&projects).Error; err != nil {
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

	if err := config.DB.Preload("Comments").Preload("Likes").First(&project, id).Error; err != nil {
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

	file, err := c.FormFile("image")
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't Save Image"})
			return
		}

		img, _ := imaging.Open(path)
		resize := imaging.Resize(img, 800, 0, imaging.Lanczos)
		if err := imaging.Save(resize, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image can't be resized"})
			return
		}
		project.Image = path
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating Project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func PutProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project Not found"})
		return
	}

	var input models.ProjectUpdateInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Format"})
		return
	}

	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.Description != nil {
		updates["description"] = *input.Description
	}
	// images
	file, err := c.FormFile("image")
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't Save Image"})
			return
		}

		img, _ := imaging.Open(path)
		resize := imaging.Resize(img, 800, 0, imaging.Lanczos)
		if err := imaging.Save(resize, path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image can't be resized"})
			return
		}

		if project.Image != "" {
			_ = os.Remove(project.Image)
		}

		updates["image"] = path
	}

	// skills
	if input.Skills != nil {
		updates["skills"] = datatypes.JSONSlice[string](*input.Skills)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update Requested"})
		return
	}

	if err := config.DB.Model(&project).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating project"})
		return
	}

	c.JSON(http.StatusOK, project)

}

func DeleteProject(c *gin.Context) {
	var project models.Project

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	// Check if id format is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Check if the project exists
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project Not found"})
		return
	}

	// remove project
	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project succesfully deleted"})
}

func LikeProject(c *gin.Context) {
	var project models.Project
	var user models.User

	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Id"})
		return
	}

	// Preload first in order to make sure the relationship are loaded.
	if err := config.DB.Preload("Likes").First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not Found"})
		return
	}

	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Id unavailable"})
		return
	}

	userIDint, ok := userId.(int)

	if err := config.DB.First(&user, uint(userIDint)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No User Found"})
		return
	}

	liked := false

	for _, u := range project.Likes {
		if u.ID == user.ID {
			liked = true
			break
		}
	}
	if liked {
		if err := config.DB.Model(&project).Association("Likes").Delete(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't remove Like"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Like Removed"})
	} else {
		if err := config.DB.Model(&project).Association("Likes").Append(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't Add Like"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Like Added"})
	}

}
