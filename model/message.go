package model

type Message struct {
	Id         int64  `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true"` //消息id
	ToUserId   int64  `json:"to_user_id,omitempty"`                              // 该消息接收者的id
	FromUserId int64  `json:"from_user_id,omitempty"`                            // 该消息发送者的id
	Content    string `json:"content,omitempty"`                                 // 消息内容
	CreateTime int64  `json:"create_time,omitempty"`                             // 消息创建时间
}

func (Message) TableName() string {
	return "messages"
}
