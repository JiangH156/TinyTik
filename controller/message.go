package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/repository"
	"TinyTik/resp"
	"TinyTik/utils/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	resp.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	redis := common.GetRedisClient()
	if user, exist := redis.UserLoginInfo(token); exist { //用户存在，生成聊天key
		userIdB, _ := strconv.Atoi(toUserId)
		curMessage := model.Message{
			Content:    content,
			CreateTime: int64(time.Now().Unix()),
			ToUserId:   int64(userIdB),
			FromUserId: user.Id,
		}
		err := repository.SendMsg(curMessage)
		if err != nil {
			// 处理错误，例如记录日志或返回错误响应
			logger.Error(err) // 记录错误日志

			// 返回错误响应
			c.JSON(http.StatusInternalServerError, resp.Response{
				StatusCode: 1,
				StatusMsg:  "Failed to send message",
			})
			return
		}
		fmt.Println("发送数据成功")
		c.JSON(http.StatusOK, resp.Response{StatusCode: 0, StatusMsg: "send success"})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	preMsgTime, _ := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)

	//TODO
	fmt.Print("进到chat里面了")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		msgList, _ := repository.GetMeassageList(user.Id, int64(userIdB), preMsgTime)
		fmt.Printf("msgList: %v\n", msgList)
		c.JSON(http.StatusOK, ChatResponse{Response: resp.Response{StatusCode: 0, StatusMsg: "pull success"}, MessageList: msgList})
	} else {
		c.JSON(http.StatusOK, resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
