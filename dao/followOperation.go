package dao

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Follow struct {
	Id         int       `gorm:"id"`
	BeFollow   int64     `gorm:"be_follow"`
	Follow     int64     `gorm:"follow"`
	IsDel      int       `gorm:"is_del"`
	UpdateTime time.Time `gorm:"update_time"`
}

func (Follow) TableName() string {
	return "follow"
}

type FollowStatus struct {
	Id    int `gorm:"id"`
	IsDel int `gorm:"is_del"`
}

type FollowActionDao struct {
}

var followActionDao *FollowActionDao
var followActionOnce sync.Once

func NewFollowActionDaoInstance() *FollowActionDao {
	followActionOnce.Do(
		func() {
			followActionDao = &FollowActionDao{}
		})
	return followActionDao
}

// GetFollowState 返回是否关注了 若返回结果为空，说明没有数据 如果返回结果不为空通过IsDel判断关系是否已经被删除了
func (*FollowActionDao) GetFollowState(followId int64, beFollow int64) (*FollowStatus, error) {
	var IsFollow FollowStatus
	resultIsFollow := db.Table("follow").Select("id, is_del").Where("follow = ? and be_follow = ?", followId, beFollow).Limit(1).Find(&IsFollow)
	if resultIsFollow.RowsAffected == 0 || resultIsFollow.Error != nil {
		if resultIsFollow.RowsAffected == 0 || errors.Is(resultIsFollow.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, resultIsFollow.Error
		}
	}
	return &IsFollow, nil
}

// UpdateStatus 更新数据，更新是否删除关系或恢复关系
func (*FollowActionDao) UpdateStatus(followId int64, beFollowId int64, id int, isDel int) error {
	follow := Follow{
		Id:         id,
		BeFollow:   beFollowId,
		Follow:     followId,
		IsDel:      isDel,
		UpdateTime: time.Now(),
	}
	result := db.Model(follow).Save(&follow)
	return result.Error
}

// CreateFollow 创建关系
func (*FollowActionDao) CreateFollow(followId int64, beFollowId int64) error {
	follow := Follow{
		BeFollow:   beFollowId,
		Follow:     followId,
		IsDel:      0,
		UpdateTime: time.Now(),
	}
	result := db.Model(follow).Select("be_follow", "follow", "is_del", "update_time").Create(&follow)
	return result.Error
}
