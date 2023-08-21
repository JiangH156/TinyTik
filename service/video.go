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
	model.Video
	model.User    `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"`
	IsFavorite    bool  `json:"is_favorite"`
}

type FeedService interface {
	Feed(c context.Context, latestTime time.Time, userId int64) (*[]VideoList, time.Time, error)
	Publish(c context.Context, video *model.Video) error
	PublishList(c context.Context, userId int64) (*[]VideoList, error)
	GetRespVideo(ctx context.Context, videoList *[]model.Video, userId int64) (*[]VideoList, error)
}

var _ FeedService = (*VideoList)(nil)

func initVideo() {
	feedService = &VideoList{}
}
func NewVideo() *VideoList {
	once.Do(initVideo)
	return feedService
}

func (v *VideoList) Feed(c context.Context, latestTime time.Time, userId int64) (*[]VideoList, time.Time, error) {

	videos, earliestTime, err := repository.NewFeed().GetVideosByLatestTime(c, latestTime)
	if err != nil {
		return nil, time.Now(), err
	}

	respVideo, err := v.GetRespVideo(c, videos, userId)
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
	respVideo, err := v.GetRespVideo(c, videos, userId)
	if err != nil {
		return nil, err
	}

	return respVideo, nil

}

func (v *VideoList) GetRespVideo(ctx context.Context, videoList *[]model.Video, userId int64) (*[]VideoList, error) {
	var resp []VideoList

	for _, v := range *videoList {

		var respVideo VideoList

		respVideo.Video = v

		wg := sync.WaitGroup{}
		wg.Add(5)

		////注意要加错误处理和redis
		go func(v *VideoList) {
			defer wg.Done()
			userInfo, err := repository.NewUserRepository().GetUserById(v.Video.AuthorId) //GetUserInfoByAuthorId

			if err != nil {
				//日志
				return
			}

			userInfo.Signature = "try"
			userInfo.Avatar = "http://8.130.16.80:8080/public/1.jpg"
			userInfo.BackgroundImage = "http://8.130.16.80:8080/public/3.jpg"

			v.User = userInfo

		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done()
			favoriteCount, err := repository.NewLikes().GetLikeCountByVideoId(ctx, v.Video.Id)
			if err != nil {
				//日志
				return
			}
			v.FavoriteCount = favoriteCount

		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done()

			commentCount, err := repository.NewCommentRepository().GetCommentCountByVideoId(ctx, v.Video.Id)
			if err != nil {
				//日志
				return

			}
			v.CommentCount = commentCount

		}(&respVideo)
		go func(v *VideoList) {
			defer wg.Done() //用户不存在时是默认值false
			isFavorite, err := repository.NewLikes().GetIslike(ctx, v.Video.Id, userId)
			if err != nil {
				//日志
				return
			}
			v.IsFavorite = isFavorite

		}(&respVideo)

		// 维护is_Follow字段
		go func(v *VideoList) {
			defer wg.Done()
			// 用户未登录状态
			if userId == 0 {
				return
			}
			repo := repository.GetRelaRepo()
			rel, err := repo.GetRelationById(userId, v.Video.AuthorId)
			if err != nil { // 不存在relation记录或出错
				return
			}
			if rel.Status == model.FOLLOW { // 当前登录用户已关注视频发布者用户
				respVideo.IsFollow = true
			}
		}(&respVideo)
		wg.Wait()

		resp = append(resp, respVideo)

	}
	return &resp, nil
}
