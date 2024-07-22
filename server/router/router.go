package router

import (
	"github.com/axitdhola/zipfile-insights/server/handlers"
	"github.com/gin-gonic/gin"
)

func InitRouter(userHandler handlers.UserHandler) *gin.Engine {
	r := gin.Default()

	// User routes
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", userHandler.GetUser)
	}

	return r
}
