package model

type Message struct {
	Id         int64  `json:"id,omitempty" gorm:"primary_key;autoIncrement:true"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
