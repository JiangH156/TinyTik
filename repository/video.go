package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/utils/logger"
	"context"
	"time"

	"gorm.io/gorm"
)

type VideoRepositoy interface {
	Save(ctx context.Context, video *model.Video) error

	GetVideosByUserID(ctx context.Context, userId int64) (*[]model.Video, error)
	GetVideosByLatestTime(ctx context.Context, latestTime time.Time) (*[]model.Video, time.Time, error)
	GetVideoListByLikeIdList(ctx context.Context, likeList []int64) (*[]model.Video, error)

	GetCommentCountByVideoId(ctx context.Context, videoId int64) (int64, error)
}

type videos struct {
	db *gorm.DB
}

var _ VideoRepositoy = (*videos)(nil)

func NewFeed() *videos {
	return &videos{
		db: common.GetDB(),
	}
}

func (v *videos) Save(ctx context.Context, video *model.Video) error {

	videor := video
	return v.db.Save(&videor).Error
}

func (v *videos) GetVideosByUserID(ctx context.Context, userId int64) (*[]model.Video, error) {
	var videos []model.Video
	err := v.db.Where("author_id=?", userId).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &videos, err
}

func (v *videos) GetVideosByLatestTime(ctx context.Context, latestTime time.Time) (*[]model.Video, time.Time, error) {
	var videos []model.Video
	// err := v.db.Model(&model.Video{}).Where("created_at < ?", latestTime).Order("created_at desc").Limit(30).Find(&videos).Error
	err := v.db.Model(&model.Video{}).Order("created_at desc").Limit(30).Find(&videos).Error

	if err != nil {
		return nil, time.Now(), err

	} else {
		if len(videos) == 0 {
			logger.Debug("no videos")
			return &videos, time.Now(), nil
		} else {
			return &videos, videos[len(videos)-1].CreatedAt, nil

		}
	}

}

func (v *videos) GetVideoListByLikeIdList(ctx context.Context, likeList []int64) (*[]model.Video, error) {
	var videoList []model.Video
	err := v.db.Where("id in ?", likeList).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return &videoList, nil

}
func (v *videos) GetCommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	var commentCount int64
	err := v.db.Model(&model.Comment{}).Where("video_id=?", videoId).Count(&commentCount).Error
	if err != nil {
		return -1, err
	}
	// Check if the commentCount is zero before using it
	if commentCount == 0 {
		return 0, nil
	}

	return commentCount, nil
}
