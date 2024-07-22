package main

import (
	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/db"
	"github.com/axitdhola/zipfile-insights/server/handlers"
	"github.com/axitdhola/zipfile-insights/server/router"
	"github.com/axitdhola/zipfile-insights/server/services"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}
	userDAO := dao.NewUserDao(dbConn.GetDB())

	// Initialize Services
	userService := services.NewUserService(userDAO)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Router
	r := router.InitRouter(userHandler)

	// Run the server
	r.Run(":8080")
}
