package main

import (
	"fmt"
	"log"

	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/db"
	"github.com/axitdhola/zipfile-insights/server/handlers"
	"github.com/axitdhola/zipfile-insights/server/router"
	"github.com/axitdhola/zipfile-insights/server/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully", dbConn.GetDB().Driver())
	userDAO := dao.NewUserDao(dbConn.GetDB())
	fileDAO := dao.NewFileDao(dbConn.GetDB())

	userService := services.NewUserService(userDAO)
	fileService := services.NewFileService(fileDAO)

	userHandler := handlers.NewUserHandler(userService)
	fileHandler := handlers.NewFileHandler(fileService)

	r := router.InitRouter(userHandler, fileHandler)

	r.Run(":8080")
}
