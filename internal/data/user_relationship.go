package data

import "context"

func (r *userRepo) GetUserFansListByPage(ctx context.Context, userid uint32, page uint32) (idList []uint32, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepo) GetUserFollowListByPage(ctx context.Context, userid uint32, page uint32) (isList []uint32, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepo) AddUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepo) DeleteUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepo) GetUserFollowInfo(ctx context.Context, userid uint32) (followCount, fansCount uint32, err error) {
	//TODO implement me
	panic("implement me")
}
