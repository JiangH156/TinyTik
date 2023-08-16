package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var MsgDB *gorm.DB = common.GetDB()
var messagesLock sync.RWMutex

// 发送信息

func SendMsg(msg model.Message) error {
	messagesLock.Lock()
	defer messagesLock.Unlock()

	err := MsgDB.Create(&msg).Error
	if err != nil {
		return err
	}

	return nil
}

func GetMeassageList(userIdA, userIdB, preMsgTime int64) ([]model.Message, error) {
	messagesLock.RLock()
	defer messagesLock.RUnlock()

	var messages []model.Message
	// 查询
	result := MsgDB.Table("messages").Where("((to_user_id = ? and from_user_id = ?) or (to_user_id = ? and from_user_id = ?)) AND create_time > ?",
		userIdA, userIdB, userIdB, userIdA, preMsgTime).Order("create_time ASC").Find(&messages)
	if result.Error != nil {
		// 处理查询错误
		return nil, result.Error
	}
	fmt.Printf("messages: %v\n", messages)
	return messages, nil
}
