package repository

import (
	"TinyTik/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB
var targetIpPort string
var messagesLock sync.RWMutex

func init() {
	InitMessage()
}

func InitMessage() {
	fmt.Print("执行init函数")
	//用viper读取message.yaml配置文件
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.SetConfigName("application_dev")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Printf("viper.Get(\"dsn\"): %v\n", viper.Get("datasource.dsn"))
	targetIpPort = viper.GetString("datasource.ipport")
	//连接到数据库dsn
	dsn := viper.GetString("datasource.dsn_no_db")
	fmt.Println("dns:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //在 GORM v2 中，数据库连接是由 GORM 管理的连接池自动管理的，并且不需要手动关闭连接
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("??????????")
	// 创建数据库
	err = db.Exec("CREATE DATABASE IF NOT EXISTS TinyTik").Error
	if err != nil {
		log.Fatal(err)
	}

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold: time.Second, // 慢 SQL 阈值
	// 		LogLevel:      logger.Info, // Log level
	// 		Colorful:      true,        // 禁用彩色打印
	// 	},
	// )
	// 连接到 TinyTik 数据库
	dsn = viper.GetString("datasource.dsn")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("成功连接到数据库！")

	// 创建数据表
	err = DB.Table("messages").AutoMigrate(&model.Message{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("数据迁移成功！")
}

// 发送信息
func SendMsg(msg model.Message) {
	conn, err := net.Dial("tcp", targetIpPort)
	if err != nil {
		fmt.Printf("无法连接到目标地址: %v\n", err)
		return
	}
	defer conn.Close()

	messagesLock.Lock()
	defer messagesLock.Unlock()
	DB.Create(&msg)

	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Message转换失败: %v\n", err)
		return
	}
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Printf("发送消息失败: %v\n", err)
		return
	}
	fmt.Println("消息发送成功！")
}

func GetMeassageList(userIdA, userIdB, preMsgTime int64) ([]model.Message, error) {
	messagesLock.RLock()
	defer messagesLock.RUnlock()

	var messages []model.Message
	// 查询
	result := DB.Table("messages").Where("((to_user_id = ? and from_user_id = ?) or (to_user_id = ? and from_user_id = ?)) AND create_time > ?",
		userIdA, userIdB, userIdB, userIdA, preMsgTime).Order("create_time ASC").Find(&messages)
	if result.Error != nil {
		// 处理查询错误
		return nil, result.Error
	}
	fmt.Printf("messages: %v\n", messages)
	return messages, nil
}
