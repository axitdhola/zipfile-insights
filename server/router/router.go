package router

import (
	"time"

	"github.com/axitdhola/zipfile-insights/server/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(userHandler handlers.UserHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	userGroup := r.Group("/users")
	{
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
	}

	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", userHandler.GetUser)
		fileGroup.GET("/:uid", userHandler.GetUser)
	}

	return r
}
