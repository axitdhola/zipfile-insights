package handlers

import (
	"net/http"

	"github.com/axitdhola/zipfile-insights/server/services"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUser(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userServices services.UserService) UserHandler {
	return &userHandler{userService: userServices}
}

func (u *userHandler) GetUser(c *gin.Context) {
	user := u.userService.GetUser(1)
	c.JSON(http.StatusOK, user)
}
