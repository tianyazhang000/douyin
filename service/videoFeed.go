package service

import (
	"douyin/dao"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type VideoFeed struct {
	NextTime   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}

// Video
type Video struct {
	Author        User   `json:"author"`         // 视频作者信息
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	ID            int64  `json:"id"`             // 视频唯一标识
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`       // 视频播放地址
	Title         string `json:"title"`          // 视频标题
}

// 视频作者信息
//
// User
type User struct {
	FollowCount   int64  `json:"follow_count"`   // 关注总数
	FollowerCount int64  `json:"follower_count"` // 粉丝总数
	ID            int64  `json:"id"`             // 用户id
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
	Name          string `json:"name"`           // 用户名称
}
type VideoService struct {
}

func (*VideoService) GetVideoFeed(latestTime int64, userId int64) *VideoFeed {
	var videoFeed VideoFeed
	//查询视频列表
	videos, _ := videoListDao.VideoList(latestTime)
	//查询视频的信息
	videoList, code := GetVideoInformation(videos, userId)
	videoFeed.StatusCode = code
	videoFeed.VideoList = videoList
	var msg string
	if userId == -1 {
		msg = "用户未登录"
	} else {
		msg = "正常"
	}
	videoFeed.StatusMsg = &msg
	var nextTime int64
	if len(videos) <= 0 {
		nextTime = time.Now().Unix()
	} else {
		nextTime = videos[len(videos)-1].CreateDate.Unix()
	}

	videoFeed.NextTime = &nextTime
	return &videoFeed
}

// GetVideoInformation 通过videoListDao.VideoList查到的的列表，得到作者、评论数、喜欢数等信息
func GetVideoInformation(videos []dao.Video, userId int64) ([]Video, int64) {
	var waitGroup sync.WaitGroup
	videoFeedStaticCode := int64(0)
	videoList := make([]Video, len(videos))
	waitGroup.Add(len(videos))
	for index, indexVideo := range videos {
		//多线程查询video的信息们
		indexVideo := indexVideo
		//此处需要对indexVideo进行阴影处理，否则会产生一个数据竞争
		go func(video *dao.Video, userId int64, index int) {
			defer waitGroup.Done()
			author, err1 := videoListDao.AuthorInformation(video.AuthorId, userId)
			if err1 != nil {
				if errors.Is(err1, gorm.ErrRecordNotFound) {
					author.UserName = "用户已注销"
				} else {
					videoFeedStaticCode++
				}
			}
			countCount, err2 := videoListDao.VideoCommentCount(video.Id)
			if err2 != nil {
				if errors.Is(err2, gorm.ErrRecordNotFound) {
					countCount = 0
				} else {
					videoFeedStaticCode++
				}
			}
			IsFavorite, err3 := videoListDao.FavoriteStatus(video.Id, userId)
			if err3 != nil {
				if errors.Is(err3, gorm.ErrRecordNotFound) {
					IsFavorite = false
				} else {
					videoFeedStaticCode++
				}
			}
			favoriteCount, err4 := videoListDao.VideoFavoriteCount(video.Id)
			if err4 != nil {
				if errors.Is(err4, gorm.ErrRecordNotFound) {
					favoriteCount = 0
				} else {
					videoFeedStaticCode++
				}
			}
			videoList[index] = Video{
				Author: User{
					FollowCount:   author.FollowCount,
					FollowerCount: author.FollowerCount,
					ID:            author.Id,
					IsFollow:      author.IsFollow,
					Name:          author.UserName,
				},
				CommentCount:  int64(countCount),
				CoverURL:      dao.ResourceDirectory + video.CoverUrl,
				FavoriteCount: int64(favoriteCount),
				ID:            video.Id,
				IsFavorite:    IsFavorite,
				PlayURL:       dao.ResourceDirectory + video.PlayUrl,
				Title:         video.Title,
			}
		}(&indexVideo, userId, index)
	}
	waitGroup.Wait()
	return videoList, videoFeedStaticCode
}
