package main

import (
	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/handlers"
	"github.com/axitdhola/zipfile-insights/server/router"
	"github.com/axitdhola/zipfile-insights/server/services"
)

func main() {
	// Initialize DAOs
	userDAO := dao.NewUserDao()

	// Initialize Services
	userService := services.NewUserService(userDAO)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Router
	r := router.InitRouter(userHandler)

	// Run the server
	r.Run(":8080")
}
