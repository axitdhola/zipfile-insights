package dao

import (
	"database/sql"

	"github.com/axitdhola/zipfile-insights/server/models"
)

type UserDao interface {
	GetUser(id int) (models.User, error)
	CreateUser(user models.User) (models.User, error)
}

type userDaoImpl struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDaoImpl{
		db: db,
	}
}

func (u *userDaoImpl) GetUser(id int) (models.User, error) {
	var user models.User
	res, err := u.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	for res.Next() {
		err = res.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (u *userDaoImpl) CreateUser(user models.User) (models.User, error) {
	_, err := u.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}