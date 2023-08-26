package service

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/resp"
	"strconv"

	"gorm.io/gorm"
)

type CommentService struct {
	DB *gorm.DB
}

func NewCommentService() *CommentService {
	return &CommentService{
		DB: common.GetDB(),
	}
}

// 保存评论
func (c *CommentService) SaveComment(comment *model.Comment) (int64, error) {
	commentRepo := repository.NewCommentRepository()
	tx := c.DB.Begin()

	commentID, err := commentRepo.CreateComment(tx, comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return commentID, nil
}

// 删除评论
func (c *CommentService) DeleteComment(commentID, videoID string) error {
	commentRepo := repository.NewCommentRepository()
	tx := c.DB.Begin()

	if err := commentRepo.DeleteCommentById(tx, commentID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 获取视频的评论列表
func (c *CommentService) GetCommentList(videoIdStr string) ([]resp.CommentResponse, error) {

	videoID, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	commentRepo := repository.NewCommentRepository()
	comments, err := commentRepo.GetCommentsByVideoID(videoID)
	if err != nil {
		return nil, err
	}

	// 获取评论中的 user_id 列表
	userIDs := make([]int64, len(comments))
	for i, comment := range comments {
		userIDs[i] = comment.User
	}

	// 查询用户列表
	users, err := commentRepo.FindUsersByUserIDs(userIDs)
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

func findUserByID(users []*model.User, userID int64) *model.User {
	for _, user := range users {
		if user.Id == userID {
			return user
		}
	}
	return &model.User{}
}
