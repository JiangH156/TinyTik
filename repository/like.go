package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/utils/logger"
	"context"

	"gorm.io/gorm"
)

type LikeRepositoy interface {
	FavoriteAction(ctx context.Context, tx *gorm.DB, userId int64, videoId int64, Liked bool) error
	GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error)
	GetLikeCountByVideoId(ctx context.Context, videoId int64) (int64, error)
	GetIslike(ctx context.Context, videoId int64, userId int64) (bool, error)
}

type likes struct {
	db *gorm.DB
}

var _ LikeRepositoy = (*likes)(nil)

func NewLikes() *likes {
	return &likes{
		db: common.GetDB(),
	}
}

func (l *likes) FavoriteAction(ctx context.Context, tx *gorm.DB, userId int64, videoId int64, liked bool) error {

	var id int64
	like := model.Like{UserId: userId, VideoId: videoId, Liked: liked}

	err := tx.Model(&model.Like{}).Select("id").Where("user_id = ? and video_id =?", userId, videoId).Find(&id).Error
	if err != nil {
		return err
	}

	if id != 0 {
		like.Id = id
	}

	err = tx.Save(&like).Error
	if err != nil {
		return err
	}
	return nil

}

func (l *likes) GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var likeList []int64
	err := l.db.Model(&model.Like{}).Select("video_id").Where("user_id = ? and liked = ?", userId, true).Find(&likeList).Error
	if err != nil {
		logger.Debug("func (l *likes) GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {")
		return nil, err
	}
	logger.Debug("likeList", likeList)

	return likeList, nil
}

func (l *likes) GetLikeCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	var likeCount int64
	err := l.db.Model(&model.Like{}).Where("video_id=? and liked =?", videoId, true).Count(&likeCount).Error
	if err != nil {
		return -1, err
	}
	return likeCount, nil
}

func (l *likes) GetIslike(ctx context.Context, videoId int64, userId int64) (bool, error) {

	var isLike bool
	err := l.db.Model(&model.Like{}).Select("liked").Where("video_id=? and user_id=?", videoId, userId).Find(&isLike).Error

	if err != nil {
		logger.Debug("repository.GetIslike false")
		return false, err
	}
	return isLike, nil
}
