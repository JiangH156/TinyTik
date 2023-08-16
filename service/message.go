package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"

	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func NewMessageService() *MessageService {
	return &MessageService{
		DB: common.GetDB(),
	}
}

func (m *MessageService) SendMsg(msg *model.Message) error {
	messageRepo := repository.NewMessageRepository()
	tx := m.DB.Begin()

	if err := messageRepo.CreateMessage(tx, msg); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *MessageService) GetMeassageList(userIdA, userIdB, preMsgTime int64) ([]model.Message, error) {
	messageRepo := repository.NewMessageRepository()
	tx := m.DB.Begin()

	messages, err := messageRepo.GetMessages(tx, userIdA, userIdB, preMsgTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return messages, nil
}
