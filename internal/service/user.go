package service

import (
	"context"

	pb "user/proto/api/user/v1"
)

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserInfo, error) {
	userInfo, err := s.uc.CreateUser(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:        userInfo.ID,
		Username:  userInfo.Name,
		Bio:       userInfo.Bio,
		AvatarUrl: userInfo.AvatarUrl,
		Tags:      userInfo.Tags,
	}, nil
}

func (s *UserService) GetUserByAccountID(ctx context.Context, req *pb.GetUserInfoByAccountIDRequest) (*pb.UserInfo, error) {
	userInfo, err := s.uc.GetUserByAccountID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Username:  userInfo.Name,
		Bio:       "",
		AvatarUrl: userInfo.AvatarUrl,
		Tags:      userInfo.Tags,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {
	return &pb.UserInfo{}, nil
}
func (s *UserService) GetUserSimpleInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserSimpleInfo, error) {
	return &pb.UserSimpleInfo{}, nil
}
func (s *UserService) GetMultipleUsersSimpleInfo(ctx context.Context, req *pb.GetMultipleUsersSimpleInfoRequest) (*pb.GetMultipleUsersSimpleInfoReply, error) {
	return &pb.GetMultipleUsersSimpleInfoReply{}, nil
}
func (s *UserService) SaveUser(ctx context.Context, req *pb.SaveUserRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *UserService) GetTagList(ctx context.Context, req *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	return &pb.GetTagListReply{}, nil
}
func (s *UserService) GetUploadAvatarURL(ctx context.Context, req *pb.GetUploadAvatarURLRequest) (*pb.GetUploadAvatarURLReply, error) {
	return &pb.GetUploadAvatarURLReply{}, nil
}
