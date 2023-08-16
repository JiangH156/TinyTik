package service

import (
	"TinyTik/model"
	"TinyTik/repository"

	"context"
	"sync"
	"time"
)

var (
	once        sync.Once
	feedService *VideoList
)

// FavoriteList  在service的响应结构体
type VideoList struct {
	VideoS        model.Video
	UserS         model.User `json:"author"`
	FavoriteCount int64      `json:"favorite_count"`
	CommentCount  int64      `json:"comment_count"`
	IsFavorite    bool       `json:"is_favorite"`
}

type FeedService interface {
	Feed(c context.Context, latestTime time.Time) (*[]VideoList, time.Time, error)
	Publish(c context.Context, video *model.Video) error
	PublishList(c context.Context, userId int64) (*[]VideoList, error)
	GetRespVideo(ctx context.Context, videoList *[]model.Video) (*[]VideoList, error)
}

var _ FeedService = (*VideoList)(nil)

func initVideo() {
	feedService = &VideoList{}
}
func NewVideo() *VideoList {
	once.Do(initVideo)
	return feedService
}

func (v *VideoList) Feed(c context.Context, latestTime time.Time) (*[]VideoList, time.Time, error) {

	videos, earliestTime, err := repository.NewFeed().GetVideosByLatestTime(c, latestTime)
	if err != nil {
		return nil, time.Now(), err
	}

	respVideo, err := v.GetRespVideo(c, videos)
	if err != nil {
		return nil, time.Now(), err
	}

	return respVideo, earliestTime, nil

}

func (v *VideoList) Publish(c context.Context, video *model.Video) error {
	err := repository.NewFeed().Save(c, video)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoList) PublishList(c context.Context, userId int64) (*[]VideoList, error) {
	videos, err := repository.NewFeed().GetVideosByUserID(c, userId)
	if err != nil {
		return nil, err
	}
	respVideo, err := v.GetRespVideo(c, videos)
	if err != nil {
		return nil, err
	}

	return respVideo, nil

}

func (v *VideoList) GetRespVideo(ctx context.Context, videoList *[]model.Video) (*[]VideoList, error) {
	var resp []VideoList

	for _, v := range *videoList {

		var respVideo VideoList

		respVideo.VideoS = v

		wg := sync.WaitGroup{}
		wg.Add(4)

		////注意要加错误处理和redis
		go func(v *VideoList) {
			defer wg.Done()
			userInfo, err := repository.NewUserRepository().GetUserById(v.VideoS.AuthorId) //GetUserInfoByAuthorId
			userInfo.Signature = "try"
			userInfo.Avatar = "http://localhost:8080/public/1.jpg"
			userInfo.BackgroundImage = "http://localhost:8080/public/1.jpg"
			if err != nil {
				//日志
				return
			}
			v.UserS = userInfo

		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done()
			favoriteCount, err := repository.NewLikes().GetLikeCountByVideoId(ctx, v.VideoS.Id)
			if err != nil {
				//日志
				return
			}
			v.FavoriteCount = favoriteCount

		}(&respVideo)
		go func(v *VideoList) {
			defer wg.Done()
			commentRepo := repository.NewCommentRepository()
			commentCount, err := commentRepo.GetCommentCountByVideoId(ctx, v.VideoS.Id)
			if err != nil {
				//日志
				return

			}
			v.CommentCount = commentCount

		}(&respVideo)
		go func(v *VideoList) {
			defer wg.Done()
			isFavorite, err := repository.NewLikes().GetIslike(ctx, v.VideoS.Id, v.VideoS.AuthorId)
			if err != nil {
				//日志
				return

			}
			v.IsFavorite = isFavorite

		}(&respVideo)
		wg.Wait()

		resp = append(resp, respVideo)

	}
	return &resp, nil
}
