package model

type User struct {
	Id              int64  `gorm:"id;primary_key;autoIncrement;comment:用户id"`
	Name            string `gorm:"name;type:varchar(32);omitempty;comment:用户名称"`
	FollowCount     int64  `gorm:"follow_count;omitempty;comment:关注总数"`
	FollowerCount   int64  `gorm:"follower_count;omitempty;comment:粉丝总数"`
	IsFollow        bool   `gorm:"is_follow;omitempty;comment:是否关注"` //true-已关注，false-未关注
	Avatar          string `gorm:"avatar;omitempty;comment:用户头像"`
	BackgroundImage string `gorm:"background_image;omitempty;comment:用户个人页顶部大图"`
	Signature       string `gorm:"signature;omitempty;comment:个人简介"`
	TotalFavorited  int64  `gorm:"total_favorited;omitempty;comment:获赞数量"`
	WorkCount       int64  `gorm:"work_count;omitempty;comment:作品数量"`
	FavoriteCount   int64  `gorm:"favorite_count;omitempty;comment:点赞数量"`
}
