package main

import (
	"douyin/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	//设置静态目录为common文件夹
	r.Static("/static", "./common")
	r.Static("/favicon.ico", "/common/tinicon.ico.png")
	apiRouter := r.Group("/douyin")

	//// basic apis
	//apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/feed", controller.VideoFeed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	//关注列表
	apiRouter.GET("/relation/follow/list/", controller.FollowerList)
	//粉丝列表
	apiRouter.GET("/relation/follower/list/", controller.BeFollowerList)
	apiRouter.GET("/")
}
