package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/utils/logger"
	"fmt"
	"strconv"

	"context"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
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

	///
	GetAuthorInfoByredis(c context.Context, userId int64, authorId int64) (*model.User, error)
	GetAuthorInfoBymysql(c context.Context, userId int64, authorId int64) (*model.User, error)

	GetLikeCountByRedis(c context.Context, videoId int64) (int64, error)
	GetCommentCountByRedis(c context.Context, videoId int64) (int64, error)

	GetIslikeByRedis(c context.Context, videoId int64, userId int64) (bool, error)
	GetIsFollowByRedis(c context.Context, userId int64, authorId int64) (bool, error)
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

	db := common.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		err := repository.NewFeed().Save(c, tx, video)
		if err != nil {
			return err
		}
		if err := repository.NewUserRepository().UpdateUserInfo(tx, video.AuthorId, 0, 3); err != nil {
			return err
		}
		return nil
	})
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
		wg.Add(4)
		go func(v *VideoList) {
			defer wg.Done()
			userInfo, err := v.GetAuthorInfoByredis(ctx, userId, v.Video.AuthorId)
			if err != nil {
				logger.Debug("获取用户信息错误", err)
				return
			}
			v.User = *userInfo
		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done()
			favoriteCount, err := v.GetLikeCountByRedis(ctx, v.Video.Id)
			if err != nil {
				//日志
				logger.Debug("获取喜欢数目错误", err)
				return
			}
			v.FavoriteCount = favoriteCount

		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done()

			commentCount, err := v.GetCommentCountByRedis(ctx, v.Video.Id)
			if err != nil {
				//日志
				logger.Debug("获取评论数量错误", err)
				return

			}
			v.CommentCount = commentCount

		}(&respVideo)

		go func(v *VideoList) {
			defer wg.Done() //用户不存在时是默认值false
			// 用户未登录状态
			if userId == 0 {
				logger.Debug("用户未登录状态")
				return
			}
			isFavorite, err := v.GetIslikeByRedis(ctx, v.Video.Id, userId)

			if err != nil {
				//日志
				logger.Debug("获取是否喜欢错误", err)
				return
			}
			v.IsFavorite = isFavorite

		}(&respVideo)

		wg.Wait()

		resp = append(resp, respVideo)

	}
	return &resp, nil
}

func (v *VideoList) GetAuthorInfoByredis(c context.Context, userId int64, authorId int64) (*model.User, error) {
	redisData, err := common.RedisA.HGetAll(c, fmt.Sprintf("authorId:%v", strconv.FormatInt(authorId, 10))).Result()
	if err == redis.Nil || len(redisData) == 0 {
		userInfo, err := v.GetAuthorInfoBymysql(c, userId, authorId)
		if err != nil {
			logger.Debug("v.GetUserInfoBymysql 获取用户信息错误")
			return nil, err
		}

		err = common.RedisA.HSet(c,
			fmt.Sprintf("authorId:%v", strconv.FormatInt(authorId, 10)),
			"id", strconv.FormatInt(userInfo.Id, 10),
			"name", userInfo.Name,
			"follow_count", strconv.FormatInt(userInfo.FollowCount, 10),
			"follower_count", strconv.FormatInt(userInfo.FollowerCount, 10),
			"avatar", userInfo.Avatar,
			"background_image", userInfo.BackgroundImage,
			"signature", userInfo.Signature,
			"total_favorited", strconv.FormatInt(userInfo.TotalFavorited, 10),
			"work_count", strconv.FormatInt(userInfo.WorkCount, 10),
			"favorite_count", strconv.FormatInt(userInfo.FavoriteCount, 10),
		).Err()
		if err != nil {
			logger.Debug("common.RedisA.HMSet  mysql设置redis用户信息错误", err)
			return nil, err
		}
		err = common.RedisA.Expire(c,
			fmt.Sprintf("authorId:%v", strconv.FormatInt(authorId, 10)),
			10*time.Minute).Err()
		if err != nil {
			logger.Debug("common.RedisA.Expire  redis过期时间错误", err)
			return nil, err
		}

		return userInfo, nil

	} else if err != nil {

		logger.Debug(err)
		return nil, err
	} else {

		var author *model.User = &model.User{}

		author.Id, err = strconv.ParseInt(redisData["id"], 10, 64)
		if err != nil {
			logger.Debug(err)
			return nil, err
		}
		author.FollowCount, _ = strconv.ParseInt(redisData["follow_count"], 10, 64)
		author.FollowerCount, _ = strconv.ParseInt(redisData["follower_count"], 10, 64)
		author.TotalFavorited, _ = strconv.ParseInt(redisData["total_favorited"], 10, 64)
		author.WorkCount, _ = strconv.ParseInt(redisData["work_count"], 10, 64)
		author.FavoriteCount, _ = strconv.ParseInt(redisData["favorite_count"], 10, 64)

		// 字符串字段直接赋值
		author.Name = redisData["name"]
		author.Avatar = redisData["avatar"]
		author.BackgroundImage = redisData["background_image"]
		author.Signature = redisData["signature"]

		if userId != 0 {
			// 自己视频
			if userId != authorId {
				isFollow, err := v.GetIsFollowByRedis(c, userId, authorId)

				if err != nil {
					//日志
					logger.Debug("获取是否关注错误", err)
					return nil, err
				}
				author.IsFollow = isFollow

			}

		}

		return author, nil

	}

}

func (v *VideoList) GetAuthorInfoBymysql(c context.Context, userId int64, authorId int64) (*model.User, error) {

	userInfo, err := repository.NewUserRepository().GetUserById(authorId)
	if err != nil {
		logger.Debug("repository.NewUserRepository().GetUserById 获取用户信息错误")
		return nil, err
	}

	return &userInfo, nil
}

func (v *VideoList) GetLikeCountByRedis(c context.Context, videoId int64) (int64, error) {
	redisData, err := common.RedisA.Get(c, fmt.Sprintf("likeCount:%d", videoId)).Result()
	if err == redis.Nil || len(redisData) == 0 {
		likeCount, err := repository.NewLikes().GetLikeCountByVideoId(c, videoId)
		if err != nil {
			logger.Debug(err)
			return 0, err
		}
		err = common.RedisA.Set(c, fmt.Sprintf("likeCount:%d", videoId), strconv.FormatInt(likeCount, 10), 10*time.Minute).Err()
		if err != nil {
			logger.Debug(err)
			return 0, err
		}
		return likeCount, nil
	} else if err != nil {
		logger.Debug(err)
		return 0, err
	}
	count, err := strconv.ParseInt(redisData, 10, 64)
	if err != nil {
		logger.Debug(err)
		return 0, err
	}

	return count, nil

}

func (v *VideoList) GetCommentCountByRedis(c context.Context, videoId int64) (int64, error) {
	redisData, err := common.RedisA.Get(c, fmt.Sprintf("commentCount:%d", videoId)).Result()
	if err == redis.Nil || len(redisData) == 0 {
		commentCount, err := repository.NewCommentRepository().GetCommentCountByVideoId(c, videoId)
		if err != nil { // 不存在relation记录或出错
			logger.Debug(err)
			return 0, err
		}
		err = common.RedisA.Set(c, fmt.Sprintf("commentCount:%d", videoId), strconv.FormatInt(commentCount, 10), 10*time.Minute).Err()
		if err != nil {
			logger.Debug(err)
			return 0, err
		}
		return commentCount, nil
	} else if err != nil {
		logger.Debug(err)
		return 0, err
	}
	count, err := strconv.ParseInt(redisData, 10, 64)
	if err != nil {
		logger.Debug(err)
		return 0, err
	}

	return count, nil

}
func (v *VideoList) GetIslikeByRedis(c context.Context, videoId int64, userId int64) (bool, error) {

	redsiData, err := common.RedisA.Get(c, fmt.Sprintf("isLike:%d:%d", videoId, userId)).Result()
	if err == redis.Nil || len(redsiData) == 0 {

		isLike, err := repository.NewLikes().GetIslike(c, videoId, userId)
		if err != nil {
			logger.Debug("GetRelationById  获取关注信息错误")
			return false, err
		}

		//更新redis
		err = common.RedisA.Set(c, fmt.Sprintf("isLike:%d:%d", videoId, userId), strconv.FormatBool(isLike), 10*time.Minute).Err()
		if err != nil {
			logger.Debug(err)
			return false, err
		}
		return isLike, nil

	} else if err != nil {

		logger.Debug(err)
		return false, err
	}

	like, err := strconv.ParseBool(redsiData)
	if err != nil {
		logger.Debug(err)
		return false, err
	}

	return like, nil
}

func (v *VideoList) GetIsFollowByRedis(c context.Context, userId int64, authorId int64) (bool, error) {

	redisData, err := common.RedisA.Get(c, fmt.Sprintf("isFollow:%v:%v", userId, authorId)).Result()
	if err == redis.Nil || len(redisData) == 0 {
		rel, err := repository.GetRelaRepo().GetRelationById(userId, authorId)
		if err != nil { // 不存在relation记录或出错
			logger.Debug("GetRelationById  获取关注信息错误")
			return false, err
		}

		var status bool
		if rel.Status == model.FOLLOW {
			status = true
		}

		//更新redis
		err = common.RedisA.Set(c, fmt.Sprintf("isFollow:%v:%v", userId, authorId), status, 10*time.Minute).Err()

		if err != nil {
			logger.Debug(err)
			return status, err
		}

		return status, nil

	} else if err != nil {

		logger.Debug(err)
		return false, err
	}

	isFollow, err := strconv.ParseBool(redisData)
	if err != nil {
		logger.Debug(err)
		return false, err
	}

	return isFollow, nil

}
