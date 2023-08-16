package middleware

import (
	"TinyTik/common"
	"TinyTik/resp"
	"TinyTik/utils/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中的认证信息是否存在
		//authToken := c.GetHeader("Authorization")
		authToken := ""
		if c.Request.Method == "GET" {
			authToken = c.Query("token")
		} else if c.Request.Method == "POST" {
			authToken = c.PostForm("token")
		}

		if authToken == "" {
			fmt.Println("未提供有效的身份验证信息")
			//logger.Error("未提供有效的身份验证信息")
			resp.Resp(c, http.StatusUnauthorized, &resp.Response{
				StatusCode: 401,
				StatusMsg:  "未提供有效的身份验证信息",
			})
			// 跳转登录界面,前端没有处理这个情况
			//c.Redirect(http.StatusFound, "/douyin/user/login/")
			return
		}
		// 检查令牌有效性
		redis := common.GetRedisClient()
		if !redis.TokenIsExist(authToken) {
			logger.Error("token不存在，验证失败")
			resp.Resp(c, http.StatusUnauthorized, &resp.Response{
				StatusCode: 401,
				StatusMsg:  "token验证失败",
			})
			// 跳转登录界面
			//c.Redirect(http.StatusFound, "/douyin/user/login/")
			return
		}
		c.Next()
	}
}
