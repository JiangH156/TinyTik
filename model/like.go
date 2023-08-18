package model

type Like struct {
	Id      int64 `json:"id"`
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
	Liked   bool  `json:"liked"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

func (l *Like) TableName() string {
	return "likes"
}
