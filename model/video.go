package model

import "time"

type Video struct {
	Id        int64     `json:"id"`
	AuthorId  int64     `json:"author_id"`
	Title     string    `json:"title"`
	PlayUrl   string    `json:"play_url"`
	CoverUrl  string    `json:"cover_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (v *Video) TableName() string {
	return "video"
}

// FavoriteList  在service的响应结构体
// type VideoList struct {
// 	VideoS        Video
// 	UserS         User  `json:"author"`
// 	FavoriteCount int64 `json:"favorite_count"`
// 	CommentCount  int64 `json:"comment_count"`
// 	IsFavorite    bool  `json:"is_favorite"`
// }
