package controller

import (
	"TinyTik/resp"
	"TinyTik/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	resp.Response
	VideoList []service.VideoList `json:"video_list,omitempty"`
	NextTime  int64               `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	// userId := c.Query("user_id")
	s := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeParse, _ := time.ParseInLocation(s, latestTime, loc)

	feedS := service.NewVideo()
	feedVideo, earliestTime, err := feedS.Feed(c, timeParse)

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
