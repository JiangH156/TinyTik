package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	resp.Response
	CommentList []resp.CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	resp.Response
	resp.CommentResponse `json:"comment,omitempty"`
}

var commentIdSequence = int64(0) //commentId的id号

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoIdInt, _ := strconv.Atoi(videoIdStr)
	redis := common.GetRedisClient()
	if user, exist := redis.UserLoginInfo(token); exist { //需要一个根据token找到user的接口
		if actionType == "1" { //发送评论

			//更新redis中的commentCount
			err := common.RedisA.Incr(c, fmt.Sprintf("commentCount:%v", videoIdStr)).Err()
			if err != nil {
				logger.Debug(err)
				return
			}
			text := c.Query("comment_text")
			tempComment := model.Comment{
				User:       int64(user.Id),
				Content:    text,
				CreateDate: time.Now().Format("05-01"),
				VideoId:    int64(videoIdInt),
			}
			CommentService := service.NewCommentService()
			commentID, err1 := CommentService.SaveComment(&tempComment)
			if err1 != nil {
				logger.Debug(err1)
				return
			}
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: resp.Response{StatusCode: 0},
				CommentResponse: resp.CommentResponse{
					Id:         commentID,
					User:       user,
					Content:    text,
					CreateDate: time.Now().Format("05-01")},
			})

			return
		} else if actionType == "2" { //删除评论

			//更新redis中的commentCount
			err := common.RedisA.Decr(c, fmt.Sprintf("commentCount:%v", videoIdStr)).Err()
			if err != nil {
				logger.Debug(err)
				return
			}
			comment_id := c.Query("comment_id")
			video_id := c.Query("video_id")
			CommentService := service.NewCommentService()
			CommentService.DeleteComment(comment_id, video_id)
		}
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")

	//获取评论
	CommentService := service.NewCommentService()
	commentList, err := CommentService.GetCommentList(video_id)
	if err != nil {
		logger.Fatal(err)
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    resp.Response{StatusCode: -1},
			CommentList: []resp.CommentResponse{},
		})
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    resp.Response{StatusCode: 0},
			CommentList: commentList,
		})
	}
}
