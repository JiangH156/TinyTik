package repository

import (
	"TinyTik/common"
	"TinyTik/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		DB: common.GetDB(),
	}
}

func (messageRepo *MessageRepository) CreateMessage(tx *gorm.DB, msg *model.Message) error {
	if err := tx.Create(msg).Error; err != nil {
		return err
	}
	return nil
}

func (messageRepo *MessageRepository) GetMessages(tx *gorm.DB, userIdA, userIdB, preMsgTime int64) ([]model.Message, error) {
	var messages []model.Message
	// 查询
	result := tx.Table("messages").Where("((to_user_id = ? and from_user_id = ?) or (to_user_id = ? and from_user_id = ?)) AND create_time > ?",
		userIdA, userIdB, userIdB, userIdA, preMsgTime).Order("create_time ASC").Find(&messages)
	if result.Error != nil {
		// 处理查询错误
		return nil, result.Error
	}
	return messages, nil
}
