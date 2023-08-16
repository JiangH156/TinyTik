package repository

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/resp"
	"strconv"
	"sync"

	"TinyTik/utils/logger"

	"gorm.io/gorm"
)

var CommentDB *gorm.DB = common.GetDB()
var commentsLock sync.RWMutex

// 保存评论
func SaveComment(comment *model.Comment) error {
	commentsLock.Lock()
	defer commentsLock.Unlock()
	result := CommentDB.Table("comments").Create(comment)
	if result.Error != nil {
		return result.Error
	}
	logger.Info("评论已成功保存，ID为：", comment.Id)
	return nil
}

// 删除评论
func DeleteComment(commentID, videoID string) error {
	// 开始事务
	tx := CommentDB.Begin()

	// 加锁
	commentsLock.Lock()
	defer commentsLock.Unlock()

	// 删除评论
	var comment model.Comment
	result := tx.Table("comments").Delete(&comment, commentID)
	if result.Error != nil {
		// 回滚事务
		tx.Rollback()
		return result.Error
	}
	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		// 处理提交事务错误
		tx.Rollback()
		return err
	}

	logger.Info("记录已成功删除: %v", comment)

	// // 通过 videoID 找到对应的视频，将视频评论总数 commentCount 减一：commentCount--
	// video, err := GetVideoByID(videoID) // 需要一个根据 videoID 获取视频的函数
	// if err != nil {
	// 	return err
	// }

	// video.CommentCount--
	// err = UpdateVideo(video) // 需要一个更新视频信息的函数
	// if err != nil {
	// 	return err
	// }
	return nil
}

var videoLock sync.RWMutex //确保在视频评论总数更新期间不会发生竞争条件或并发冲突
func GetVideoByID(videoID string) (*model.Video, error) {
	// 这是一个空函数，需要根据实际情况进行实现
	return nil, nil
}

func UpdateVideo(video *model.Video) error {
	// 这是一个空函数，需要根据实际情况进行实现
	return nil
}

// 获取评论列表，这个评论列表需要包含user这个对象一起返回
func GetCommentList(videoIdStr string) ([]resp.CommentResponse, error) {
	videoIdInt, err := strconv.Atoi(videoIdStr)
	if err != nil {
		return nil, err
	}
	videoID := int64(videoIdInt)

	commentsLock.RLock()

	var comments []*model.Comment
	result := CommentDB.Table("videos").Where("video_id = ?", videoID).Find(&comments)
	if result.Error != nil {
		commentsLock.RUnlock()
		return nil, result.Error
	}
	commentsLock.RUnlock()

	// 获取评论中的 user_id 列表
	userIDs := make([]int64, len(comments))
	for i, comment := range comments {
		userIDs[i] = comment.User
	}

	// 查询用户列表
	users, err := FindUsersByUserIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// 构建 CommentResponse 列表
	commentsResponse := make([]resp.CommentResponse, len(comments))
	for i, comment := range comments {
		user := findUserByID(users, comment.User)
		commentResponse := resp.CommentResponse{
			Id:         comment.Id,
			User:       *user,
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		}
		commentsResponse[i] = commentResponse
	}

	return commentsResponse, nil
}

func FindUsersByUserIDs(userIDs []int64) ([]*model.User, error) {
	//commentsLock.Lock()
	//defer commentsLock.Unlock()   避免死锁的发生这里的查询不加锁
	// 查询用户
	var users []*model.User
	result := CommentDB.Table("users").Where("id IN ?", userIDs).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func findUserByID(users []*model.User, userID int64) *model.User {
	for _, user := range users {
		if user.Id == userID {
			return user
		}
	}
	return nil
}
