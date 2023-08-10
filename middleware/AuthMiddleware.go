package middleware

import (
	"TinyTik/common"
	"TinyTik/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中的认证信息是否存在
		//authToken := c.GetHeader("Authorization")
		authToken := c.Query("token")
		if authToken == "" {
			fmt.Println("未提供有效的身份验证信息")
			//logger.Error("未提供有效的身份验证信息")
			resp.Resp(c, http.StatusUnauthorized, &resp.Response{
				StatusCode: 401,
				StatusMsg:  "未提供有效的身份验证信息",
			})
			// 跳转登录界面
			c.Redirect(http.StatusFound, "/douyin/user/login/")
			return
		}
		// 检查令牌有效性
		if !ValidateToken(authToken) {
			//logger.Error("token验证失败")
			resp.Resp(c, http.StatusUnauthorized, &resp.Response{
				StatusCode: 401,
				StatusMsg:  "token验证失败",
			})
			// 跳转登录界面
			c.Redirect(http.StatusFound, "/douyin/user/login/")
			return
		}
		c.Next()
	}
}

// 验证token有效性
func ValidateToken(vToken string) bool {
	// redis中是否缓存token
	_, err := common.Redis.GetToken(vToken)
	return err != nil
}
