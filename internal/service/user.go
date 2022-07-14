package service

import (
	"context"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"user/internal/biz"
	pb "user/proto/api/user/v1"
)

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserInfo, error) {
	user, err := s.uc.CreateUser(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Name,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

func (s *UserService) GetUserInfoByAccountID(ctx context.Context, req *pb.GetUserInfoByAccountIDRequest) (*pb.UserInfo, error) {
	user, err := s.uc.GetUserByAccountID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Name,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {
	user, err := s.uc.GetUserInfoByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Name,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
		Tags:      user.Tags,
	}, nil
}

func (s *UserService) GetMultipleUsersInfo(ctx context.Context, req *pb.GetMultipleUsersInfoRequest) (*pb.GetMultipleUsersInfoReply, error) {
	users, err := s.uc.GetMultipleUsersInfo(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	usersInfo := make([]*pb.UserInfo, len(users))
	for i, user := range users {
		usersInfo[i] = &pb.UserInfo{
			Id:        user.ID,
			Username:  user.Name,
			Bio:       user.Bio,
			AvatarUrl: user.AvatarUrl,
			Tags:      user.Tags,
		}
	}

	return &pb.GetMultipleUsersInfoReply{Info: usersInfo}, nil
}

func (s *UserService) UserEdit(ctx context.Context, req *pb.UserEditRequest) (*empty.Empty, error) {
	_, err := s.uc.SaveUserInfo(ctx, &biz.UserInfo{
		ID:        req.Id,
		Name:      *req.Username,
		Bio:       *req.Bio,
		AvatarUrl: *req.AvatarUrl,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *UserService) GetTagList(ctx context.Context, req *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	tagMap, err := s.uc.GetUserTagList(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTagListReply{
		Tags: tagMap,
	}, nil
}

func (s *UserService) GetUploadAvatarURL(ctx context.Context, req *pb.GetUploadAvatarURLRequest) (*pb.GetUploadAvatarURLReply, error) {
	url, err := s.uc.GetUploadURL(ctx, req.Id, req.Crc32, req.Sha1)
	if err != nil {
		return nil, err
	}

	return &pb.GetUploadAvatarURLReply{
		Url: url,
	}, nil
}

func (s *UserService) BanUser(ctx context.Context, req *pb.BanUserRequest) (*empty.Empty, error) {
	_, err := s.uc.BanUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *UserService) DeBanUser(ctx context.Context, req *pb.DeBanUserRequest) (*empty.Empty, error) {
	_, err := s.uc.DeBanUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *UserService) UserEditTags(ctx context.Context, req *pb.UserEditTagsRequest) (*empty.Empty, error) {
	_, err := s.uc.EditUserTags(ctx, req.Id, req.Tags)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
