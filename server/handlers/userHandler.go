package handlers

import (
	"net/http"

	"github.com/axitdhola/zipfile-insights/server/models"
	"github.com/axitdhola/zipfile-insights/server/services"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUser(c *gin.Context)
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
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

func (u *userHandler) RegisterUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)
	user, err := u.userService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *userHandler) LoginUser(c *gin.Context) {
	user := models.User{}
	c.BindJSON(&user)
	user, err := u.userService.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
