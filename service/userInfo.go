package service

import "douyin/dao"

func GetUserInfo(userId int64) (dao.Author, error) {

	return NewGetUserInfoFlow(userId).Do()
}

type GetUserInfoFlow struct {
	userId int64
	User   *dao.Author
}

func NewGetUserInfoFlow(userId int64) *GetUserInfoFlow {
	return &GetUserInfoFlow{userId: userId}
}

func (f *GetUserInfoFlow) Do() (dao.Author, error) {
	if err := f.CheckUserId(); err != nil {
		return dao.Author{}, err
	}
	if err := f.GetUserInfo(); err != nil {
		return dao.Author{}, err
	}
	return *f.User, nil
}
func (f *GetUserInfoFlow) CheckUserId() error {
	if _, err := dao.FindUserById(f.userId); err != nil {
		return err
	}
	return nil
}
func (f *GetUserInfoFlow) GetUserInfo() error {
	var err error
	if f.User, err = dao.NewVideoDaoInstance().GetUserInformation(f.userId); err != nil {
		return err
	}
	return nil
}
