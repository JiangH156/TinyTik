package controller

import (
	"TinyTik/common"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AllFavoriteList struct {
	Res    resp.Response
	Videos *[]service.VideoList `json:"video_list"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	// userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	var userId int64
	token := c.Query("token")
	if token == "" {
		logger.Debug("tokennnnnnnnnnnnnnnnnn")
	}
	redis := common.GetRedisClient()
	if user, exist := redis.UserLoginInfo(token); exist {
		userId = user.Id
	} else {
		logger.Debug("user not exist")
	}

	actionTypeInt32 := c.Query("action_type")
	if actionTypeInt32 == "" {
		// 处理空字符串的情况
		logger.Debug("空字符串")
	}

	actionType, err := strconv.ParseInt(actionTypeInt32, 10, 64)
	if err != nil {
		// 处理解析错误
		logger.Debug("转换错误")
	}

	likeSerVice := service.NewlikeSerVice()

	if err := likeSerVice.FavoriteAction(c, userId, videoId, actionType); err != nil {
		c.JSON(http.StatusInternalServerError, resp.Response{
			StatusCode: -1,
			StatusMsg:  "FavoriteAction false",
		})
	} else {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 0,
			StatusMsg:  "FavoriteAction success",
		})
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	videoService := service.NewlikeSerVice()

	videoList, err := videoService.FavoriteList(c, userId)
	if err != nil {
		logger.Debug("videoService.FavoriteList")
		c.JSON(http.StatusInternalServerError, AllFavoriteList{
			Res: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Like list false",
			},
			Videos: nil})

	} else {

		c.JSON(http.StatusOK, AllFavoriteList{
			Res: resp.Response{
				StatusCode: 0,
				StatusMsg:  "Like list success",
			},
			Videos: videoList,
		})
	}
}
