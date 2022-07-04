package service

import (
	"github.com/google/wire"
	"user/internal/biz"
	v1 "user/proto/api/user/v1"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserService,
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
