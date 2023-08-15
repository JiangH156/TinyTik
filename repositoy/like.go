package repositoy

import (
	"TinyTik/common"
	"TinyTik/model"
	"context"

	"gorm.io/gorm"
)

type LikeRepositoy interface {
	FavoriteAction(ctx context.Context, userId int64, videoId int64, Liked bool) error
	GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error)
	GetLikeCountByVideoId(ctx context.Context, videoId int64) (int64, error)
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

func (l *likes) FavoriteAction(ctx context.Context, userId int64, videoId int64, liked bool) error {

	var like model.Like
	like = model.Like{UserId: userId, VideoId: videoId, Liked: liked}
	if err := l.db.Save(&like).Error; err != nil {

		return err

	}
	return nil
}
func (l *likes) GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var likeList []int64
	err := l.db.Select("video_id").Where("user_id=? and liked=?", userId, true).Find(&likeList).Error
	if err != nil {
		return nil, err
	}

	return likeList, nil
}
func (l *likes) GetLikeCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	var likeCount int64
	err := l.db.Model(&model.Like{}).Where("video_id=? and like =?", videoId, true).Count(&likeCount).Error
	if err != nil {
		return -1, err
	}
	return likeCount, nil

}