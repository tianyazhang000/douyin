package service

import (
	"douyin/dao"
)

type GetFollowerListFlow struct {
	TokenUserId  int64
	UserId       int64        `json:"user_id"`
	FollowerList []dao.Author `json:"user_list"`
}

func GetFollowerList(tokenUserId int64, userId int64) ([]dao.Author, error) {
	return NewGetFollowerList(tokenUserId, userId).Do()
}

func NewGetFollowerList(tokenUserId int64, userId int64) *GetFollowerListFlow {
	return &GetFollowerListFlow{
		TokenUserId: tokenUserId,
		UserId:      userId,
	}
}

func (f *GetFollowerListFlow) Do() ([]dao.Author, error) {
	if err := f.CheckUserId(); err != nil {
		return []dao.Author{}, err
	}
	if err := f.GetList(); err != nil {
		return []dao.Author{}, err
	}
	return f.FollowerList, nil
}

// 检查用户Id是否存在
func (f *GetFollowerListFlow) CheckUserId() error {
	if _, err := dao.FindUserById(f.UserId); err != nil {
		return err
	}
	return nil
}

// 先获取Id的粉丝列表，再根据列表返回粉丝信息
func (f *GetFollowerListFlow) GetList() error {
	followIdList := dao.NewFollowListDaoInstance().GetFollowIdList(f.UserId)
	for _, followId := range followIdList {
		var user *dao.Author
		var err error
		if user, err = dao.NewVideoDaoInstance().AuthorInformation(followId, f.TokenUserId); err != nil {
			return err
		}
		f.FollowerList = append(f.FollowerList, *user)
	}
	return nil
}
