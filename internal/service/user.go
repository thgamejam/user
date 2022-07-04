package service

import (
	"context"
	pb "user/proto/api/user/v1"
)

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	userInfo, err := s.uc.CreateUser(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		UserInfo: &pb.UserInfo{
			Name:      userInfo.Name,
			AvatarUrl: userInfo.AvatarUrl,
			WorkCount: userInfo.WorkCount,
			FansCount: userInfo.FansCount,
			Tags:      userInfo.Tags,
		},
	}, nil
}

func (s *UserService) GetUserByAccountID(ctx context.Context, req *pb.GetUserByAccountIDRequest) (*pb.GetUserByAccountIDReply, error) {
	userInfo, err := s.uc.GetUserByAccountID(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByAccountIDReply{
		UserInfo: &pb.UserInfo{
			Name:      userInfo.Name,
			AvatarUrl: userInfo.AvatarUrl,
			WorkCount: userInfo.WorkCount,
			FansCount: userInfo.FansCount,
			Tags:      userInfo.Tags,
		},
	}, nil
}
