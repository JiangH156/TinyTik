package model

type Video struct {
	Id            int64  `gorm:"id;primary_key;autoIncrement;comment:视频唯一标识"`
	Author        User   `gorm:"author;comment:视频作者信息"`
	PlayUrl       string `gorm:"play_url,omitempty;comment:视频播放地址"`
	CoverUrl      string `gorm:"cover_url,omitempty;comment:视频封面地址"`
	FavoriteCount int64  `gorm:"favorite_count,omitempty;comment:视频的点赞总数"`
	CommentCount  int64  `gorm:"comment_count,omitempty;comment:视频的评论总数"`
	IsFavorite    bool   `gorm:"is_favorite,omitempty;comment:是否点赞"` // true-已点赞，false-未点赞
	Title         string `gorm:"title;omitempty;comment:视频标题"`
}
