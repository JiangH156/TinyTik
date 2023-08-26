package controller

import (
	"TinyTik/common"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	resp.Response
	VideoList []service.VideoList `json:"video_list"`
	NextTime  int64               `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
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

	latestTime := c.Query("latest_time")

	// logger.Debug(latestTime, " url时间")

	latestTimeUnix, err := strconv.ParseInt(latestTime, 10, 64)
	if err != nil {
		logger.Debug("Error parsing UNIX timestamp:", err)
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeParse := time.Unix(latestTimeUnix, 0).In(loc)

	if timeParse.After(time.Now()) {
		timeParse = time.Now()
	}

	feedS := service.NewVideo()
	feedVideo, earliestTime, err := feedS.Feed(c, timeParse, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, FeedResponse{
			Response:  resp.Response{StatusCode: -1, StatusMsg: "Feed False"},
			VideoList: *feedVideo,
			NextTime:  earliestTime.Unix(),
		})

	} else {

		c.JSON(http.StatusOK, FeedResponse{
			Response:  resp.Response{StatusCode: 0, StatusMsg: "Feed OK"},
			VideoList: *feedVideo,
			NextTime:  earliestTime.Unix(),
		})
	}

}
