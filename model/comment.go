package model

type Comment struct {
	Id         int64  `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true"`
	User       int64  `json:"user_id"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64  `json:"vedeo_id,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentResponse struct {
	Id         int64  `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
