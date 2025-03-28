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

// Register 用户注册
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

// Login 用户登录
func (s *UserService) Login(user *model.User) (error, string) {
	//先判断用户是否存在
	existingUser, err := s.UserRepository.GetUserByUsername(user.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err, ""
	}
	if existingUser == nil {
		return myerrors.ErrUserNotExisted, ""
	}

	//判断密码是否正确
	if !auth.CheckPasswordHash(user.Password, existingUser.Password) {
		return myerrors.ErrPasswordWrong, ""
	} else {
		//生成并返回token
		token, err := auth.GenerateToken(existingUser.Username, existingUser.ID)
		if err != nil {
			return err, ""
		}
		return nil, token

	}

}
