package controller

import (
	"TinyTik/resp"
	"TinyTik/service"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	resp.Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// latestTime := c.Query("latest_time")
	// userID := c.Query("user_id")
	// s := "2023-07-13 22:58:44"
	// loc, _ := time.LoadLocation("Asia/Shanghai")
	// timeP, _ := time.ParseInLocation(latestTime, s, loc)

	// S := service.NewVideo()
	// sC := S.feed()

	// c.JSON(http.StatusOK, FeedResponse{
	// 	Response:  resp.Response{StatusCode: 0},
	// 	VideoList: DemoVideos,
	// 	NextTime:  time.Now().Unix(),
	// })

}
