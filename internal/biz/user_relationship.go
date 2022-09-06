package biz

import "context"

func (uc *UserUseCase) GetUserFansList(ctx context.Context, userid uint32, page uint32) (idList []uint32, err error) {
	return uc.repo.GetUserFansListByPage(ctx, userid, page)
}

func (uc *UserUseCase) GetUserFollowList(ctx context.Context, userid uint32, page uint32) (isList []uint32, err error) {
	return uc.repo.GetUserFollowListByPage(ctx, userid, page)
}

func (uc *UserUseCase) GetUserFollowInfo(ctx context.Context, userid uint32) (followCount, fansCount uint32, err error) {
	return uc.repo.GetUserFollowInfo(ctx, userid)
}

func (uc *UserUseCase) FollowUser(ctx context.Context, userid, followUserid uint32) (err error) {
	return uc.repo.AddUserFollowInfo(ctx, userid, followUserid)
}

func (uc *UserUseCase) CancelFollowUser(ctx context.Context, userid, followUserid uint32) (err error) {
	return uc.repo.DeleteUserFollowInfo(ctx, userid, followUserid)
}
