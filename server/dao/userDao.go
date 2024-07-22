package dao

import (
	"github.com/axitdhola/zipfile-insights/server/models"
)

type UserDao interface {
	GetUser(id int) models.User
}

type userDaoImpl struct{}

func NewUserDao() UserDao {
	return &userDaoImpl{}
}

func (u *userDaoImpl) GetUser(id int) models.User {
	return models.User{Id: id, Name: "AKSHIT", Email: "akshit@123", Password: "akshit123"}
}
