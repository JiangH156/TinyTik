package repository

import (
	"TinyTik/common"
	"TinyTik/model"

	"context"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		DB: common.GetDB(),
	}
}

// 保存评论
func (c *CommentRepository) CreateComment(tx *gorm.DB, comment *model.Comment) (int64, error) {
	result := tx.Table("comments").Create(comment)
	if result.Error != nil {
		return 0, result.Error
	}
	// 获取插入后的评论对象，包括自动生成的ID
	return comment.Id, nil
}

// 删除评论
func (c *CommentRepository) DeleteCommentById(tx *gorm.DB, commentID string) error {

	var comment model.Comment
	result := tx.Table("comments").Delete(&comment, commentID)
	if result.Error != nil {
		// 回滚事务
		tx.Rollback()
		return result.Error
	}
	return nil
}

// 通过视频ID查找评论
func (c *CommentRepository) GetCommentsByVideoID(videoID int64) ([]model.Comment, error) {
	var comments []model.Comment
	result := c.DB.Table("comments").Where("video_id = ?", videoID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func (c *CommentRepository) FindUsersByUserIDs(userIDs []int64) ([]*model.User, error) {
	// 查询用户
	var users []*model.User
	result := c.DB.Table("users").Where("id IN ?", userIDs).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (c *CommentRepository) GetCommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	var commentCount int64
	err := c.DB.Model(&model.Comment{}).Where("video_id=? ", videoId).Count(&commentCount).Error
	if err != nil {
		return -1, err
	}
	return commentCount, nil
}
