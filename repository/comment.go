package repository

import (
	"TinyTik/model"
	"TinyTik/resp"
	"fmt"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var CommentDB *gorm.DB
var commentsLock sync.RWMutex

func InitComment() {
	//用viper读取message.yaml配置文件
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.SetConfigName("application_dev")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	//连接到数据库dsn
	dsn := viper.GetString("mysql.dsn_no_db")
	fmt.Println("dns:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //在 GORM v2 中，数据库连接是由 GORM 管理的连接池自动管理的，并且不需要手动关闭连接
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("??????????")
	// 创建数据库
	err = db.Exec("CREATE DATABASE IF NOT EXISTS TinyTik").Error
	if err != nil {
		logger.Fatal(err)
	}

	// 连接到 TinyTik 数据库
	dsn = viper.GetString("mysql.dsn")
	CommentDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	// 创建数据表
	err = CommentDB.Table("messages").AutoMigrate(&model.Message{}, &model.User{})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("数据迁移成功！")
}

// 保存评论
func SaveComment(comment *model.Comment) error {
	commentsLock.Lock()
	defer commentsLock.Unlock()
	result := CommentDB.Table("comments").Create(comment)
	if result.Error != nil {
		return result.Error
	}
	logger.Println("评论已成功保存，ID为：", comment.Id)
	return nil
}

// 删除评论
func DeleteComment(commentID, videoID string) error {
	commentsLock.Lock()
	defer commentsLock.Unlock()

	var comment model.Comment
	result := CommentDB.Table("comments").Delete(&comment, commentID)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("记录已成功删除: %v", comment)

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
