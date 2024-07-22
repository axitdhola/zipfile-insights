package services

import (
	"errors"
	"fmt"

	"github.com/axitdhola/zipfile-insights/server/dao"
	"github.com/axitdhola/zipfile-insights/server/models"
	"github.com/axitdhola/zipfile-insights/server/util"
)

type UserService interface {
	GetUser(id int) (models.User, error)
	RegisterUser(user models.User) (models.LoginResponse, error)
	LoginUser(user models.User) (models.LoginResponse, error)
}

type userServiceImpl struct {
	userDao dao.UserDao
}

func NewUserService(userDao dao.UserDao) UserService {
	return &userServiceImpl{userDao: userDao}
}

func (u *userServiceImpl) GetUser(id int) (models.User, error) {
	if id == 0 {
		return models.User{}, errors.New("invalid user id")
	}
	return u.userDao.GetUser(id)
}

func (u *userServiceImpl) RegisterUser(user models.User) (models.LoginResponse, error) {
	fmt.Println("user", user)
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return models.LoginResponse{}, errors.New("invalid user details")
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	user.Password = hashedPassword
	res, err := u.userDao.CreateUser(user)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{Id: res.Id, Name: res.Name, Email: res.Email}, err
}

func (u *userServiceImpl) LoginUser(user models.User) (models.LoginResponse, error) {
	if user.Email == "" || user.Password == "" {
		return models.LoginResponse{}, errors.New("invalid user details")
	}
	res, err := u.userDao.GetUser(user.Id)
	if err != nil {
		return models.LoginResponse{}, err
	}

	err = util.CheckPassword(user.Password, res.Password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{Id: res.Id, Name: res.Name, Email: res.Email}, nil
}
