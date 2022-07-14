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
	// GetUserInfoByUserID 通过用户ID获取用户信息
	GetUserInfoByUserID(ctx context.Context, userID uint32) (user *UserInfo, err error)
	// GetUserInfoByUserIDList 批量获取用户信息
	GetUserInfoByUserIDList(ctx context.Context, userId []uint32) (user []UserInfo, err error)
	// GetUserTags 获取用户所有标签
	GetUserTags(ctx context.Context, userID uint32) (tags map[uint32]string, err error)
	//CreateUser 创建用户
	CreateUser(ctx context.Context, accountID uint32) (user *UserInfo, err error)
	// SaveUser 修改用户
	SaveUser(ctx context.Context, user *UserInfo) (ok bool, err error)
	// AddUserTags 添加用户标签
	AddUserTags(ctx context.Context, userId []uint32, tagID uint32) (ok bool, err error)
	// BanUser 封禁用户
	BanUser(ctx context.Context, userID uint32) (ok bool, err error)
	// DeBanUser 解封用户
	DeBanUser(ctx context.Context, userID uint32) (ok bool, err error)
}
