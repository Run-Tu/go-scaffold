package service

import (
	"github.com/Run-Tu/go-scaffold/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) GetALLUsers() ([]models.User, error) {
	var users []models.User
	us.DB.Find(&users)

	return users, nil
}
