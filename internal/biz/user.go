package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserInfo struct {
	ID        uint32   `json:"id"`         // 用户id
	AccountID uint32   `json:"account_id"` // 账户id
	Name      string   `json:"name"`       // 用户名
	Bio       string   `json:"bio"`        // 简介
	Tags      []string `json:"tags"`       // 标签
	AvatarUrl string   `json:"avatar_url"` // 头像链接
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUseCase new a User use case.
func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(logger)}
}

// GetUserByAccountID 通过账户ID获取用户信息
func (uc *UserUseCase) GetUserByAccountID(ctx context.Context, accountID uint32) (*UserInfo, error) {
	return uc.repo.GetUserByAccountID(ctx, accountID)
}

// CreateUser 根据账户ID创建用户
func (uc *UserUseCase) CreateUser(ctx context.Context, accountID uint32) (userInfo *UserInfo, err error) {
	return uc.repo.CreateUser(ctx, accountID)
}
