package service

import (
	"chat/internal/model"
	"chat/internal/repository"
	"chat/pkg/auth"
	myerrors "chat/pkg/errors"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(user *model.User) error {
	//先判断用户是否已经存在
	existingUser, err := s.UserRepository.GetUserByUsername(user.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return myerrors.ErrUserExisted
	}

	//处理明文密码
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	//存入数据库
	if err := s.UserRepository.CreateUser(user).Error; err != nil {
		return err
	}
	return nil

}
