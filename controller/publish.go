package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type VideoListResponse struct {
	Res    resp.Response
	Videos []service.VideoList `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	title := c.PostForm("title")

	videoHeader, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Response{
			StatusCode: -1,
			StatusMsg:  "Get file err",
		})
		return
	}

	// 验证 token，获取 userID
	// userID, err := verifyToken(token)
	var userId int64
	token := c.PostForm("token")
	redis := common.GetRedisClient()
	if user, exist := redis.UserLoginInfo(token); exist {
		userId = user.Id
	} else {
		logger.Debug("user not exist")
	}

	// 存储视频数据
	videoPath := fmt.Sprintf("public/%s-%s", uuid.New().String(), videoHeader.Filename)

	if err := c.SaveUploadedFile(videoHeader, videoPath); err != nil {
		c.JSON(http.StatusInternalServerError, resp.Response{
			StatusCode: -1,
			StatusMsg:  "Save file err",
		})
		return
	}
	//压缩视频
	// videoPath, err = compressVideo(videoPath)
	// if err != nil {
	// 	logger.Debug("compressVideo false:", err)
	// 	return
	// }

	urlPre := viper.GetString("server.urlPre")

	playUrl := fmt.Sprintf("%s%s", urlPre, videoPath)

	// 截取视频封面
	coverPath := generateVideoCover(videoPath)

	coverUrl := fmt.Sprintf("%s%s", urlPre, coverPath)
	logger.Debug(coverPath)

	var video model.Video
	video.AuthorId = userId
	video.CoverUrl = coverUrl
	//video.CoverUrl = "http://localhost:8080/public/3.png"
	video.CreatedAt = time.Now()
	video.PlayUrl = playUrl
	video.Title = title
	video.UpdatedAt = time.Now()

	// 存储视频信息和封面路径，这里可以将视频和封面路径存储到数据库中
	err = service.NewVideo().Publish(c, &video)

	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Response{
			StatusCode: -1,
			StatusMsg:  "Save file err",
		})

	} else {
		c.JSON(http.StatusOK, resp.Response{
			StatusCode: 0,
			StatusMsg:  videoHeader.Filename + " uploaded successfully",
		})
	}

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	videoService := service.NewVideo()
	videoList, err := videoService.PublishList(c, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, VideoListResponse{
			Res: resp.Response{
				StatusCode: -1,
				StatusMsg:  "Publish list false",
			},
			Videos: nil})

	} else {

		c.JSON(http.StatusOK, VideoListResponse{
			Res: resp.Response{
				StatusCode: 0,
				StatusMsg:  "Publish list success",
			},
			Videos: *videoList,
		})
	}
}

func generateVideoCover(videoPath string) string {
	// 使用 ffmpeg 获取视频的第一帧作为封面
	coverFilename := strings.TrimSuffix(videoPath, ".mp4") + "_cover.jpg"
	command := []string{
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		coverFilename,
	}
	cmd := exec.Command("ffmpeg", command...)
	cmd.Stderr = os.Stderr // Redirect stderr to console for error messages

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error generating cover:", err)
		return ""
	}
	return coverFilename
}

func compressVideo(inputVideoPath string) (string, error) {
	outputVideoPath := strings.TrimSuffix(inputVideoPath, ".mp4") + "CMP.mp4"
	command := []string{
		"-i", inputVideoPath,
		"-c:v", "libx264",
		//"-b:v", "1M", // 使用比特率代替 -crf
		"-crf", "18",
		"-y", // This option enables overwriting without asking
		outputVideoPath,
	}
	cmd := exec.Command("ffmpeg", command...)
	cmd.Stderr = os.Stderr // Redirect stderr to console for error messages

	err := cmd.Run()
	if err != nil {
		logger.Debug("Error compressing video:", err)
		return outputVideoPath, err
	}

	return outputVideoPath, nil
}
