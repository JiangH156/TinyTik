package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/utils/logger"
	"time"

	"context"
	"fmt"

	"gorm.io/gorm"
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

		logger.Debug("执行点赞操作")
		//在Redis中记录用户点赞状态
		err := common.RedisA.Set(ctx, fmt.Sprintf("isLike:%d:%d", videoId, userId), true, 10*time.Minute).Err()
		if err != nil {
			return err
		}
		//刷新redis中的likeCount
		err = common.RedisA.Incr(ctx, fmt.Sprintf("likeCount:%d", videoId)).Err()
		if err != nil {
			logger.Debug(err)
			return err
		}

		//使用事务
		db := common.GetDB()
		err = db.Transaction(func(tx *gorm.DB) error {
			// 在MySQL中保存用户点赞记录
			if err := likeRepositoy.FavoriteAction(ctx, tx, userId, videoId, true); err != nil {
				return err
			}
			//刷新用户信息
			if err := repository.NewUserRepository().UpdateUserInfo(tx, userId, videoId, 1); err != nil {
				return err
			}

			return nil

		})
		if err != nil {
			return err
		}

	} else if action_type == 2 { //执行取消点赞操作

		logger.Debug("执行取消点赞操作")

		// 在Redis中移除用户点赞状态
		err := common.RedisA.Del(ctx, fmt.Sprintf("isLike:%d:%d", videoId, userId)).Err()
		if err != nil {
			return err
		}
		//刷新redis中的likeCount
		err = common.RedisA.Decr(ctx, fmt.Sprintf("likeCount:%d", videoId)).Err()
		if err != nil {
			logger.Debug(err)
			return err
		}

		db := common.GetDB()
		err = db.Transaction(func(tx *gorm.DB) error {
			// 在MySQL中保存用户点赞记录
			if err := likeRepositoy.FavoriteAction(ctx, tx, userId, videoId, false); err != nil {
				return err
			}
			//刷新用户信息
			if err := repository.NewUserRepository().UpdateUserInfo(tx, userId, videoId, 2); err != nil {
				return err
			}

			return nil

		})
		if err != nil {
			return err
		}

	}
	return nil

}

// favorite_count comment_count  is_favorite       __--根据userid查找用户信息
func (l *likeSerVice) FavoriteList(ctx context.Context, userId int64) (*[]VideoList, error) {

	//错误处理
	likeIdList, err := repository.NewLikes().GetlikeIdListByUserId(ctx, userId)
	if err != nil {
		logger.Debug("GetlikeIdListByUserId")
		return nil, err
	}
	videoList, err := repository.NewFeed().GetVideoListByLikeIdList(ctx, likeIdList)
	if err != nil {
		logger.Debug("GetVideoListByLikeIdList")
		return nil, err
	}
	resp, err := feedService.GetRespVideo(ctx, videoList, userId)
	if err != nil {
		logger.Debug("GetRespVideo")
		return nil, err
	}

	return resp, nil

}
