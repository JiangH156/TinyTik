package service

import (
	"TinyTik/model"
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	once        sync.Once
	feedService *Video
)

type Video struct {
	model.Video
	Author        model.User `json:"author"`
	FavoriteCount int64      `json:"favorite_count"`
	CommentCount  int64      `json:"comment_count"`
	IsFavorite    bool       `json:"is_favorite"`
}

type FeedService interface {
	Feed(c gin.Context, latestTime time.Time, userId string) (*[]Video, time.Time, error)
	PublishList(c context.Context, userId string) (*[]Video, time.Time, error)
}

// var _ FeedService = (*Video)(nil)

func initVideo() {
	feedService = &Video{}
}
func NewVideo() *Video {
	once.Do(initVideo)
	return feedService
}

// func (v *Video) Feed(c context.Context, latestTime time.Time, userId string) (*[]Video, time.Time, error) {

// }

// func (v *Video) PublishList(c context.Context, userId string) (*[]Video, time.Time, error) {

// 	var video *[]Video
// 	G := repositoy.NewVideos()

//		videos, err := G.GetVideosByUserID(c, userId)
//	}
