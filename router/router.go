package router

import (
	"TinyTik/controller"
	"TinyTik/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.AuthMiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.AuthMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.AuthMiddleware(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.AuthMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.AuthMiddleware(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.AuthMiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.AuthMiddleware(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.AuthMiddleware(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.AuthMiddleware(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.AuthMiddleware(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", middleware.AuthMiddleware(), controller.FriendList)
	apiRouter.GET("/message/chat/", middleware.AuthMiddleware(), controller.MessageChat)
	apiRouter.POST("/message/action/", middleware.AuthMiddleware(), controller.MessageAction)
}
