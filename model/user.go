package model

import (
	"TinyTik/utils/logger"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	Id              int64  `json:"id"   gorm:"id;primary_key;autoIncrement;comment:用户id"`
	Name            string `json:"name"   gorm:"name;type:varchar(32);omitempty;comment:用户名称"`
	FollowCount     int64  `json:"follow_count" gorm:"follow_count;omitempty;comment:关注总数"`
	FollowerCount   int64  `json:"follower_count" gorm:"follower_count;omitempty;comment:粉丝总数"`
	IsFollow        bool   `json:"is_follow" gorm:"is_follow;omitempty;comment:是否关注"` //true-已关注，false-未关注
	Avatar          string `json:"avatar" gorm:"avatar;omitempty;comment:用户头像"`
	BackgroundImage string `json:"background_image" gorm:"background_image;omitempty;comment:用户个人页顶部大图"`
	Signature       string `json:"signature" gorm:"signature;omitempty;comment:个人简介"`
	TotalFavorited  int64  `json:"total_favorited" gorm:"total_favorited;omitempty;comment:获赞数量"`
	WorkCount       int64  `json:"work_count" gorm:"work_count;omitempty;comment:作品数量"`
	FavoriteCount   int64  `json:"favorite_count" gorm:"favorite_count;omitempty;comment:点赞数量"`
}

// 使用 Hook 进行合法性检查
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.FollowCount < 0 || u.FollowerCount < 0 {
		logger.Error("Invalid Update: Follow/Unfollow action yield minus value.")
		return fmt.Errorf("操作数为负")
	}
	return nil
}
