package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) CreateUser(tx *gorm.DB, user model.User) (err error) {
	if err := tx.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) GetUserByUsername(username string) (user model.User, err error) {
	if err := r.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) GetUserById(id int64) (user model.User, err error) {
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: common.GetDB(),
	}
}
