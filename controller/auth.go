package controller

import (
	"TinyTik/common"
	"TinyTik/model"
	"TinyTik/resp"
	"TinyTik/service"
	"TinyTik/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户注册 PORT /douyin/user/register/
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
	id, token, lErr := authService.Register(model.UserAuth{UserName: username, Password: password})
	if lErr.Err != nil {
		resp.Resp(c, int(lErr.HttpCode), &resp.Response{
			StatusCode: lErr.HttpCode,
			StatusMsg:  lErr.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, resp.UserLoginResponse{
		Response: resp.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

// 用户登录 POST /douyin/user/login/
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
	id, token, lErr := authService.Login(model.UserAuth{UserName: username, Password: password})
	if lErr.Err != nil {
		resp.Resp(c, int(lErr.HttpCode), &resp.Response{
			StatusCode: lErr.HttpCode,
			StatusMsg:  lErr.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, resp.UserLoginResponse{
		Response: resp.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}
