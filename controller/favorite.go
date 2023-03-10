package controller

import (
	"douyin/middle"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

var favoriteService service.FavoriteService

func FavoriteAction(context *gin.Context) {
	token := context.Query("token")
	//检查token
	userid, err := middle.ParseToken(token)
	if err != nil {
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token无效",
		})
		return
	}
	videoid := context.Query("video_id")
	actionType := context.Query("action_type")
	err = favoriteService.Update(userid, videoid, actionType)
	if err != nil {
		context.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "successfully"})
	}
}
func FavoriteList(context *gin.Context) {
	var token string
	var userId int64
	var exist bool
	var err error
	//检查token
	if token, exist = context.GetQuery("token"); !exist {
		context.JSON(http.StatusOK, FollowListResponse{
			StatusCode: -1,
			StatusMsg:  "缺少token",
		})
		return
	}
	if userId, err = middle.ParseToken(token); err != nil {
		context.JSON(http.StatusOK, FollowListResponse{
			StatusCode: -1,
			StatusMsg:  "token无效",
		})
		return
	}
	favoriteList := favoriteService.FavoriteList(userId)
	context.JSON(http.StatusOK, favoriteList)
}
