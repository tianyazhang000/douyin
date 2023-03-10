package controller

import (
	"douyin/middle"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func VideoFeed(context *gin.Context) {
	token := context.Query("token")
	var videoService service.VideoService
	//处理时间，分为携带和不携带
	latestTimeStr := context.Query("latest_time")
	if len(latestTimeStr) > 10 {
		latestTimeStr = latestTimeStr[0:10]
	}
	var latestTime int64
	if latestTimeStr == "" {
		latestTime = time.Now().Unix()
	} else {
		latestTime, _ = strconv.ParseInt(latestTimeStr, 10, 64)
	}

	//未登录
	if token == "" {
		videoFeed := videoService.GetVideoFeed(latestTime, -1)
		context.JSON(http.StatusOK, videoFeed)
		return
	} else {
		//检查token
		userId, err := middle.ParseToken(token)
		if err != nil {
			msg := "token无效"
			videoFeed := videoService.GetVideoFeed(latestTime, -1)
			videoFeed.StatusMsg = &msg
			context.JSON(http.StatusOK, videoFeed)
			return
		}
		videoFeed := videoService.GetVideoFeed(latestTime, userId)
		context.JSON(http.StatusOK, videoFeed)
	}

}
