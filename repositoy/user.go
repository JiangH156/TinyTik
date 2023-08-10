package repositoy

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

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: common.GetDB(),
	}
}
