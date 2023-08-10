package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户注册
func Register(c *gin.Context) {
	// 数据接收
	username := c.Query("username")
	password := c.Query("password")
	// 数据验证
	if err := common.ValidateUserAuth(model.UserAuth{UserName: username, Password: password}); err != nil {
		logger.Error("Data format error")
		resp.Resp(c, http.StatusBadRequest, &resp.Response{
			StatusCode: 400,
			StatusMsg:  "Data format error",
		})
		return
	}
	authService := service.NewAuthService()
	// service注册逻辑
	regUser, lErr := authService.Register(model.UserAuth{UserName: username, Password: password})
	if lErr.Err != nil {
		resp.Resp(c, int(lErr.HttpCode), &resp.Response{
			StatusCode: lErr.HttpCode,
			StatusMsg:  lErr.Msg,
		})
		return
	}
	// 生成token
	token, err := common.GenToken(username)
	if err != nil {
		logger.Error("common.GenToken error, ", err)
		resp.Resp(c, http.StatusInternalServerError, &resp.Response{
			StatusCode: 500,
			StatusMsg:  "Internal server error",
		})
		return
	}
	// 新登录的用户，直接添加到redis缓存中
	newUser := model.User{
		Id:   regUser.Id,
		Name: regUser.Name,
	}
	// json格式化user
	userBytes, err := json.Marshal(newUser)
	if err != nil {
		logger.Error("json.Marshal error, ", err)
		resp.Resp(c, http.StatusInternalServerError, &resp.Response{
			StatusCode: 500,
			StatusMsg:  "Internal server error",
		})
		return
	}

	redis := common.GetRedisClient()
	err = redis.SetUser(token, userBytes)
	if err != nil {
		logger.Error("redis.SetUser error, ", err)
		resp.Resp(c, http.StatusInternalServerError, &resp.Response{
			StatusCode: 500,
			StatusMsg:  "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, resp.UserLoginResponse{
		Response: resp.Response{StatusCode: 0},
		UserId:   regUser.Id,
		Token:    token,
	})

}

func Login(c *gin.Context) {
	// 数据接收
	username := c.Query("username")
	password := c.Query("password")
	// 数据验证
	if err := common.ValidateUserAuth(model.UserAuth{UserName: username, Password: password}); err != nil {
		logger.Error("Data format error")
		resp.Resp(c, http.StatusBadRequest, &resp.Response{
			StatusCode: 400,
			StatusMsg:  "Data format error",
		})
		return
	}
	authService := service.NewAuthService()
	// service登录逻辑
	authService.Login(model.UserAuth{UserName: username, Password: password})

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, resp.UserLoginResponse{
			Response: resp.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, resp.UserLoginResponse{
			Response: resp.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
