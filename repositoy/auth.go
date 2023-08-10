package repositoy

import (
	"TinyTik/common"
	"TinyTik/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

// 创建用户
func (auth *AuthRepository) CreateAuth(tx *gorm.DB, userAuth *model.UserAuth) error {
	if err := tx.Create(userAuth).Error; err != nil {
		return err
	}
	return nil
}

// 根据用户名获取UserAuth
func (auth *AuthRepository) GetAuthByUsername(username string) (userAuth model.UserAuth, err error) {
	if err = auth.DB.Where("user_name = ?", username).First(&model.UserAuth{}).Error; err != nil {
		return userAuth, err
	}
	return userAuth, nil
}

func (auth *AuthRepository) GetIDByUsername(userName string) (id int64, err error) {
	var userAuth model.UserAuth
	// 查询用户名对应的记录
	if err = auth.DB.Where("user_name = ?", userName).First(&userAuth).Error; err != nil {
		return 0, err
	}
	return userAuth.ID, nil
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		DB: common.GetDB(),
	}
}
