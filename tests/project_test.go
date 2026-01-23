package tests

import (
	"bytes"
	"demo/config"
	"demo/controllers"
	"demo/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Project{}, &models.Comment{})

	project := models.Project{Name: "Project Test", Description: "Description Test"}
	db.Create(&project)

	comment := models.Comment{ProjectID: project.ID, Content: "Comment Test"}
	db.Create(&comment)

	return db
}

func TestGetProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.GET("/projects", controllers.GetProjects)

	req, _ := http.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "Project Test")
	assert.Contains(t, body, "Comment Test")
}

func TestPostProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.POST("/projects", controllers.PostProject)

	project := map[string]interface{}{
		"name":        "Test Project",
		"description": "Test Description",
		"skills":      []string{"Go", "Testing", "SQLite"},
	}

	data, _ := json.Marshal(project)

	req, _ := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Contains(t, w.Body.String(), "Test Project")
	assert.Contains(t, w.Body.String(), "Test Description")
	assert.Contains(t, w.Body.String(), "Testing")
}
