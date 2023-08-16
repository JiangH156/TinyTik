package repositoy

import (
	"TinyTik/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func (message *MessageRepository) CreateMsg(tx *gorm.DB, msg *model.Message) error {
	if err := tx.Create(msg).Error; err != nil {
		return err
	}
	return nil
}
