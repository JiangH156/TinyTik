package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"errors"
	"sync"

	"gorm.io/gorm"
)

type RelaRepo struct {
	DB *gorm.DB
}

var once sync.Once
var repo *RelaRepo

func GetRelaRepo() *RelaRepo {
	// 采用单例模式
	if repo == nil {
		once.Do(func() {
			repo = &RelaRepo{}
			repo.DB = common.GetDB()
		})
	}
	return repo
}

func (r *RelaRepo) Followed(user *model.User, toUser *model.User) bool {
	rel, err := r.GetRelationById(user.Id, toUser.Id)
	if err != nil {
		return false
	}
	if rel.Status == model.FOLLOW {
		return true
	}
	return false
}

func (r *RelaRepo) GetRelationById(id int64, toId int64) (model.Relation, error) {
	rel := model.Relation{}
	res := r.DB.Where("user_id = ? AND to_user_id = ?", id, toId).First(&rel)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return rel, res.Error
	}
	return rel, nil
}

func (r *RelaRepo) GetFollowListById(id int64) ([]model.User, error) {
	users := []model.User{}
	res := r.DB.Model(&model.User{}).
		Select("users.*, IF(EXISTS(SELECT * FROM relations r1 WHERE users.id = r1.user_id AND r1.to_user_id = ? AND r1.status = ?), true, false) as is_follow", id, model.FOLLOW).
		Joins("JOIN relations r ON users.id = r.to_user_id AND r.user_id = ? AND r.status = ?", id, model.FOLLOW).
		Scan(&users)
	return users, res.Error
}

func (r *RelaRepo) GetFollowerListById(id int64) ([]model.User, error) {
	users := []model.User{}
	res := r.DB.Model(&model.User{}).
		Select("users.*, IF(EXISTS(SELECT * FROM relations r1 WHERE users.id = r1.user_id AND r1.to_user_id = ? AND r1.status = ?), true, false) as is_follow", id, model.FOLLOW).
		Joins("JOIN relations r ON users.id = r.user_id AND r.to_user_id = ? AND r.status = ?", id, model.FOLLOW).
		Scan(&users)
	return users, res.Error
}

func (r *RelaRepo) GetFriendListById(id int64) ([]model.User, error) {
	users := []model.User{}
	query := r.DB.Table("relations").Select("to_user_id").Where("user_id = ?", id)
	// 互关朋友就统一设置 is_follow 字段为 true
	res := r.DB.Model(&model.User{}).
		Select("users.*, true as is_follow").
		Joins("JOIN (?) q ON users.id = q.to_user_id", query).
		Joins("JOIN relations r ON q.to_user_id = r.user_id AND r.to_user_id = ? AND r.status = ?", id, model.FOLLOW).
		Scan(&users)
	return users, res.Error
}

func (r *RelaRepo) UpdateRelation(user model.User, toUser model.User, follow byte) error {
	// 执行事务
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Save(&user).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}
		// 更新关系表
		rel, err := r.GetRelationById(user.Id, toUser.Id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rel.UserId = user.Id
			rel.ToUserId = toUser.Id
		}
		rel.Status = follow
		if err := tx.Save(&rel).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}