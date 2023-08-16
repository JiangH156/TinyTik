package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"

	"context"
	"fmt"
	"time"
)

var (
	likeS *likeSerVice
)

type LikeSerVice interface {
	//点赞还是取消
	FavoriteAction(ctx context.Context, userId int64, videoId int64, action_type int64) error
	FavoriteList(ctx context.Context, userId int64) (*[]VideoList, error) //favorite_count comment_count  is_favorite __--根据userid查找用户信息

}

// LikeSerVice的实现结构体
type likeSerVice struct {
	likes *model.Like
}

var _ LikeSerVice = (*likeSerVice)(nil)

func initlikeS() {
	likeS = &likeSerVice{}
}

func NewlikeSerVice() *likeSerVice {
	once.Do(initlikeS)
	return likeS
}

func (l *likeSerVice) FavoriteAction(ctx context.Context, userId int64, videoId int64, action_type int64) error {
	// var ls likeSerVice
	likeRepositoy := repository.NewLikes()

	if action_type == 1 { //执行点赞操作
		// 在Redis中记录用户点赞状态
		err := common.RedisA.Set(ctx, fmt.Sprintf("%d:%d", videoId, userId), 1, 500*time.Millisecond).Err()
		if err != nil {
			return err
		}

		// 在MySQL中保存用户点赞记录

		if err := likeRepositoy.FavoriteAction(ctx, userId, videoId, true); err != nil {
			return err
		}

	} else if action_type == 2 { //执行取消点赞操作

		// 在Redis中移除用户点赞状态
		err := common.RedisA.Del(ctx, fmt.Sprintf("%d:%d", videoId, userId)).Err()
		if err != nil {
			return err
		}
		// 在MySQL中删除用户点赞记录

		if err := likeRepositoy.FavoriteAction(ctx, userId, videoId, false); err != nil {
			return err
		}

	}
	return nil

}

// type VideoList struct {
// 	VideoS        Video
// 	UserS         User
// 	FavoriteCount int64 `json:"favorite_count"`
// 	CommentCount  int64 `json:"comment_count"`
// 	IsFavorite    bool  `json:"is_favorite"`
// }

// favorite_count comment_count  is_favorite       __--根据userid查找用户信息
func (l *likeSerVice) FavoriteList(ctx context.Context, userId int64) (*[]VideoList, error) {

	//错误处理
	likeIdList, err := repository.NewLikes().GetlikeIdListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	videoList, err := repository.NewFeed().GetVideoListByLikeIdList(ctx, likeIdList)
	if err != nil {
		return nil, err
	}
	resp, err := feedService.GetRespVideo(ctx, videoList)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
