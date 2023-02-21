package controller

import (
	"douyin/middle"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PublishList(context *gin.Context) {
	token := context.Query("token")
	var videoListService service.VideoListService
	//处理传来的参数userId,并查看它是否合法
	//userId的情况交给下层处理
	if userId, error := strconv.ParseInt(context.Query("user_id"), 10, 64); error != nil {
		msg := "参数错误"
		context.JSON(http.StatusOK, service.VideoListModal{
			StatusCode: 1,
			StatusMsg:  &msg,
			VideoList:  nil,
		})
		return
	} else {
		//token为空说明用户未登录
		if token == "" {
			videoFeed := videoListService.GetVideoList(-1, userId)
			context.JSON(http.StatusOK, videoFeed)
			return
		} else {
			//用户已经登录
			//检查token
			tokenUserId, err := middle.ParseToken(token)
			if err != nil {
				//if err != nil {
				msg := "token无效"
				context.JSON(http.StatusOK, service.VideoListModal{
					StatusCode: 1,
					StatusMsg:  &msg,
					VideoList:  nil,
				})
				return
			}
			videoFeed := videoListService.GetVideoList(tokenUserId, userId)
			context.JSON(http.StatusOK, videoFeed)
		}
	}
}
