package routes

import (
	"demo/controllers"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(router *gin.Engine) {
	routesGroup := router.Group("/projects")
	{
		routesGroup.GET("/", controllers.GetProjects)
		routesGroup.GET("/:id", controllers.GetProject)
		routesGroup.POST("/", controllers.PostProject)
		routesGroup.PUT("/:id", controllers.PutProject)
	}
}
