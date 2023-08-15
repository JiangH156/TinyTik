package controller

import (
	"TinyTik/resp"
	"TinyTik/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	videoId, _ := strconv.ParseInt(c.PostForm("video_id"), 10, 64)
	userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.PostForm("action_type"), 10, 64)

	likeSerVice := service.NewlikeSerVice()

	if err := likeSerVice.FavoriteAction(c, userId, videoId, actionType); err != nil {
		c.JSON(http.StatusBadRequest, resp.Response{
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
	userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	videoService := service.NewlikeSerVice()
	videoList, err := videoService.FavoriteList(c, userId)
	if err != nil {
		c.JSON(404, resp.FavoriteList{
			Res: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Like list false",
			},
			Videos: nil})

	} else {

		c.JSON(http.StatusOK, resp.FavoriteList{
			Res: resp.Response{
				StatusCode: 0,
				StatusMsg:  "Like list success",
			},
			Videos: videoList,
		})
	}
}
