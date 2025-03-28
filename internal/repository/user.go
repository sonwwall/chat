package repository

import (
	"chat/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser 创建User
func (r *UserRepository) CreateUser(user *model.User) *gorm.DB {
	return r.db.Create(user)
}

// GetUserByUsername 通过用户名获取用户
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username=?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}
