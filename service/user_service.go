package service

import (
	"errors"

	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
	"github.com/ljcbaby/form/util"
)

type UserService struct{}

func (s *UserService) Login(user model.User) (jwt string, err error) {
	db := database.DB

	var db_user model.User
	if err := db.Where("username = ?", user.Username).First(&db_user).Error; err != nil {
		return "", err
	}

	if db_user.Password != util.MD5(user.Password+db_user.Salt) {
		return "", errors.New("PWD_MISMATCH")
	}

	jwt, err = util.GenerateToken(db_user.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *UserService) CheckUsernameExist(username string) (bool, error) {
	db := database.DB

	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *UserService) Register(user *model.User) (err error) {
	db := database.DB

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserByID(user *model.User) (err error) {
	db := database.DB

	if err := db.Where("id = ?", user.ID).First(user).Error; err != nil {
		if err.Error() == "record not found" {
			return errors.New("USER_NOT_FOUND")
		}
		return err
	}

	return nil
}
