package main

import (
	"demo/config"
	"demo/models"
	"demo/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "demo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Project Sharing
// @version 1.0
// @description Description of the Shareing Project
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	// didn't see anything on the video on that one
	router.SetTrustedProxies(nil)

	// Security middleware
	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))

	// Loading .ENV vrariables
	err := godotenv.Load()
	if err != nil {
		log.Println("file not found: .ENV")
	}

	// API router definition
	routes.ProjectRoutes(router)
	routes.UserRoutes(router)
	routes.CommentRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Connect to DB
	config.ConnectDB()

	// DB migration
	config.DB.AutoMigrate(&models.Project{}, &models.User{}, &models.Comment{})

	// Start Server on port 8000
	router.Run(":8000")
}
