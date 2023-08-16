package repositoy

import (
	"TinyTik/model"
	"TinyTik/utils/logger"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

// 保存评论
func (CommentDB *CommentRepository) CreateComment(tx *gorm.DB, comment *model.Comment) error {
	result := tx.Table("comments").Create(comment)
	if result.Error != nil {
		return result.Error
	}
	logger.Info("评论已成功保存，ID为：", comment.Id)
	return nil
}
