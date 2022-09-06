package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "user/proto/api/user/v1"
)

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserInfo, error) {
	user, err := s.uc.CreateUser(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

// GetUserInfoByAccountID 通过账户ID获取用户信息
func (s *UserService) GetUserInfoByAccountID(ctx context.Context, req *pb.GetUserInfoByAccountIDRequest) (*pb.UserInfo, error) {
	user, err := s.uc.GetUserByAccountID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

// GetUserInfo 通过用户id获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {
	user, err := s.uc.GetUserInfoByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

// GetMultipleUsersInfo 获取多个用户id
func (s *UserService) GetMultipleUsersInfo(ctx context.Context, req *pb.GetMultipleUsersInfoRequest) (*pb.GetMultipleUsersInfoReply, error) {
	users, err := s.uc.GetMultipleUsersInfo(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	usersInfo := make([]*pb.UserInfo, len(users))
	for i, user := range users {
		usersInfo[i] = &pb.UserInfo{
			Id:        user.ID,
			Username:  user.Username,
			Bio:       user.Bio,
			AvatarUrl: user.AvatarUrl,
			Tags:      user.Tags,
		}
	}

	return &pb.GetMultipleUsersInfoReply{
		Info: usersInfo,
	}, nil
}

// UserEdit 用户编辑信息
func (s *UserService) UserEdit(ctx context.Context, req *pb.UserEditRequest) (*emptypb.Empty, error) {
	// TODO 传递参数不对应，待修改

	return &emptypb.Empty{}, nil
}

// UserEditTags 修改用户tag
func (s *UserService) UserEditTags(ctx context.Context, req *pb.UserEditTagsRequest) (*emptypb.Empty, error) {
	err := s.uc.EditUserTags(ctx, req.Id, req.Tags)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetTagList 获取用户tag列表
func (s *UserService) GetTagList(ctx context.Context, req *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	tags, err := s.uc.GetUserTagList(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	tagList := make(map[uint32]string)
	for i, tag := range tags {
		tagList[uint32(i)] = tag
	}

	return &pb.GetTagListReply{
		Tags: tagList,
	}, nil
}

// GetUploadAvatarURL 获取头像上传链接
func (s *UserService) GetUploadAvatarURL(ctx context.Context, req *pb.GetUploadAvatarURLRequest) (*pb.GetUploadAvatarURLReply, error) {
	url, err := s.uc.GetUploadAvatarURL(ctx, req.Id, req.Crc32, req.Sha1)
	if err != nil {
		return nil, err
	}

	return &pb.GetUploadAvatarURLReply{
		Url: url,
	}, nil
}

// BanUser 封禁用户
func (s *UserService) BanUser(ctx context.Context, req *pb.BanUserRequest) (*emptypb.Empty, error) {
	err := s.uc.BanUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeBanUser 解封用户
func (s *UserService) DeBanUser(ctx context.Context, req *pb.DeBanUserRequest) (*emptypb.Empty, error) {
	err := s.uc.DeBanUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetUserFansList 获取用户粉丝列表
func (s *UserService) GetUserFansList(ctx context.Context, req *pb.GetUserListReq) (*pb.GetUserListReply, error) {
	useridList, err := s.uc.GetUserFansList(ctx, req.Userid, req.Page)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserListReply{
		UseridList: useridList,
	}, nil
}

// GetUserFollowList 获取用户关注列表
func (s *UserService) GetUserFollowList(ctx context.Context, req *pb.GetUserListReq) (*pb.GetUserListReply, error) {
	useridList, err := s.uc.GetUserFollowList(ctx, req.Userid, req.Page)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserListReply{
		UseridList: useridList,
	}, nil
}

// GetUserFollowInfo 获取用户关注信息
func (s *UserService) GetUserFollowInfo(ctx context.Context, req *pb.GetUserFollowInfoReq) (*pb.GetUserFollowInfoReply, error) {
	fansCount, followCount, err := s.uc.GetUserFollowInfo(ctx, req.Userid)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserFollowInfoReply{
		FollowCount: fansCount,
		FansCount:   followCount,
	}, nil
}

// FollowUser 关注用户
func (s *UserService) FollowUser(ctx context.Context, req *pb.FollowUserReq) (*emptypb.Empty, error) {
	err := s.uc.FollowUser(ctx, req.Userid, req.FollowUserId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// CancelFollowUser 取消关注
func (s *UserService) CancelFollowUser(ctx context.Context, req *pb.CancelFollowUserReq) (*emptypb.Empty, error) {
	err := s.uc.CancelFollowUser(ctx, req.Userid, req.FollowUserId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
