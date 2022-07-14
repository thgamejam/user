package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUseCase)

// UserRepo is a User repo.
type UserRepo interface {
	// GetUserByAccountID 通过账户ID获取用户
	GetUserByAccountID(ctx context.Context, accountID uint32) (user *UserInfo, err error)
	//CreateUser 创建用户
	CreateUser(ctx context.Context, accountID uint32) (user *UserInfo, err error)
}
