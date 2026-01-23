package routes

import (
	"demo/controllers"
	"demo/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	routesGroup := router.Group("/comments")

	routesGroup.Use(middlewares.Authentication())
	{
		routesGroup.POST("/", controllers.PostComment)
	}
}
