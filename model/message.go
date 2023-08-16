package model

type Message struct {
	Id         int64  `gorm:"id;primary_key;autoIncrement;comment:消息id"`
	ToUserId   int64  `gorm:"to_user_id;omitempty;comment:该消息接收者的id"`
	FromUserIf int64  `gorm:"from_user_id;omitempty;comment:该消息发送者的id"`
	Content    string `gorm:"content;omitempty;comment:消息内容"`
	CreateTime string `gorm:"create_time;omitempty;comment:消息创建时间"`
}
