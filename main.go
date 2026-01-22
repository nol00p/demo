package main

import (
	"demo/config"
	"demo/models"
	"demo/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.ProjectRoutes(router)
	routes.UserRoutes(router)

	config.ConnectDB()
	fmt.Print("Server Running on http://localhost:8000")

	config.DB.AutoMigrate(&models.Project{})
	config.DB.AutoMigrate(&models.User{})

	router.Run(":8000")
}
