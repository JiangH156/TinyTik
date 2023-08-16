package repository

import (
	"TinyTik/model"
	"TinyTik/utils/logger"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MsgDB *gorm.DB
var messagesLock sync.RWMutex

// func init() {
// 	InitMessage()
// }

func InitMessage() {
	//用viper读取message.yaml配置文件
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.SetConfigName("application_dev")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	//连接到数据库dsn
	dsn := viper.GetString("datasource.dsn_no_db")
	fmt.Println("dns:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //在 GORM v2 中，数据库连接是由 GORM 管理的连接池自动管理的，并且不需要手动关闭连接
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("??????????")
	// 创建数据库
	err = db.Exec("CREATE DATABASE IF NOT EXISTS TinyTik").Error
	if err != nil {
		logger.Fatal(err)
	}

	// 连接到 TinyTik 数据库
	dsn = viper.GetString("datasource.dsn")
	MsgDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	// 创建数据表
	err = MsgDB.Table("messages").AutoMigrate(&model.Message{})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("数据迁移成功！")
}

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
