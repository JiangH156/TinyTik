package controller

import (
	"TinyTik/model"
	"TinyTik/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	resp.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  resp.Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
