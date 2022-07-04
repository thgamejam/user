package service

import (
	"context"
	
	"user/internal/biz"
	v1 "user/proto/api/user/v1"
)

// UserService is a user service.
type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUseCase
}

// NewUserService new a user service.
func NewUserService(uc *biz.UserUseCase) *UserService {
	return &UserService{uc: uc}
}

// SayHello implements helloworld.UserServer.
func (s *UserService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateUser(ctx, &biz.User{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
