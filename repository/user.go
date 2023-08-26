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

func (r *UserRepository) UpdateUserInfo(tx *gorm.DB, userId int64, videoId int64, types int64) error {

	if types == 1 { //用户点赞加一  视频作者被点赞数+1
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		var authoId int64
		if err := tx.Model(&model.Video{}).Select("author_id").Where("id= ?", videoId).Find(&authoId).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", authoId).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error; err != nil {
			return err
		}

	} else if types == 2 { //用户点赞减一	视频作者被点赞数-1
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}

		var authoId int64
		if err := tx.Model(&model.Video{}).Select("author_id").Where("id= ?", videoId).Find(&authoId).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", authoId).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error; err != nil {
			return err
		}
	} else if types == 3 { //作品数量+1

		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: common.GetDB(),
	}
}
