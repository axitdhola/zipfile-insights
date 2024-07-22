package services

import (
	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/models"
)

type UserService interface {
	GetUser(id int) models.User
}

type userServiceImpl struct {
	userDao dao.UserDao
}

func NewUserService(userDao dao.UserDao) UserService {
	return &userServiceImpl{userDao: userDao}
}

func (u *userServiceImpl) GetUser(id int) models.User {
	return u.userDao.GetUser(id)
}
