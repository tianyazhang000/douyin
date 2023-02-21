package service

import "douyin/dao"

func GetBeFollowerListFlow(tokenUserId int64, userId int64) ([]dao.Author, error) {
	return NewGetFollowerList(tokenUserId, userId).Get()
}

func (f *GetFollowerListFlow) Get() ([]dao.Author, error) {
	if err := f.CheckUserId(); err != nil {
		return []dao.Author{}, err
	}
	if err := f.GetBeFollowList(); err != nil {
		return []dao.Author{}, err
	}
	return f.FollowerList, nil
}

// 先获取Id的关注列表，再根据列表返回关注信息
func (f *GetFollowerListFlow) GetBeFollowList() error {
	followIdList := dao.NewFollowListDaoInstance().GetBeFollowIdList(f.UserId)
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
