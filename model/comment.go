package model

type Comment struct {
	Id         int64  `gorm:"id;primary_key;autoIncrement;comment:视频评论id"`
	User       User   `gorm:"user;comment:评论用户信息"`
	Content    string `gorm:"content,omitempty;comment:评论内容"`
	CreateDate string `gorm:"create_date,omitempty;comment:评论发布日期"` //格式 mm-dd
}
