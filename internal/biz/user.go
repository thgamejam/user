package biz

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/thgamejam/pkg/util"
	"user/internal/conf"
)

type UserInfo struct {
	ID        uint32   // 用户id
	AccountID uint32   // 账户id
	Username  string   // 用户名
	Bio       string   // 简介
	Tags      []string // 标签
	AvatarUrl string   // 头像链接
}

type ModifiableUserInfo struct {
	Username util.Val[string]   // 用户名
	Bio      util.Val[string]   // 简介
	Tags     util.Val[[]uint32] // 标签
}

type UserUseCase struct {
	defaultUsernamePrefix string

	repo UserRepo
	log  *log.Helper
}

// NewUserUseCase new a User use case.
func NewUserUseCase(repo UserRepo, user *conf.User, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		defaultUsernamePrefix: user.DefaultUsernamePrefix,

		repo: repo,
		log:  log.NewHelper(logger),
	}
}

var buildDefaultUsername = func(prefix string, accountID uint32) string {
	return prefix + strconv.FormatUint(uint64(accountID), 36)
}

// GetUserByAccountID 通过账户ID获取用户信息
func (uc *UserUseCase) GetUserByAccountID(ctx context.Context, accountID uint32) (userInfo *UserInfo, err error) {
	v, err := uc.repo.GetUserInfoByAccountID(ctx, accountID)
	// TODO error
	if v.IsExist() {
		userInfo = v.Val()
	}
	return
}

// CreateUser 根据账户ID创建用户
func (uc *UserUseCase) CreateUser(ctx context.Context, accountID uint32) (userInfo *UserInfo, err error) {
	return uc.repo.CreateUser(ctx, accountID, buildDefaultUsername(uc.defaultUsernamePrefix, accountID))
}

// GetUserInfoByUserID 根据用户ID获取用户信息
func (uc *UserUseCase) GetUserInfoByUserID(ctx context.Context, userID uint32) (userInfo *UserInfo, err error) {
	v, err := uc.repo.GetUserInfoByUserID(ctx, userID)
	// TODO error
	if v.IsExist() {
		userInfo = v.Val()
	}
	return
}

// GetMultipleUsersInfo 根据用户id列表批量获取用户信息
func (uc *UserUseCase) GetMultipleUsersInfo(ctx context.Context, ids []uint32) (usersInfo []*UserInfo, err error) {
	return nil, errors.New("todo")
}

// SaveUserInfo 保存用户
func (uc *UserUseCase) SaveUserInfo(ctx context.Context, userID uint32, userInfo *ModifiableUserInfo) (err error) {
	return uc.repo.EditUserInfo(ctx, userID, userInfo)
}

// GetUserTagList 获取用户所有的tag列表
func (uc *UserUseCase) GetUserTagList(ctx context.Context, userID uint32) (tags map[uint32]string, err error) {
	return uc.repo.GetUserTags(ctx, userID)
}

// GetUploadURL 获取头像上传链接
func (uc UserUseCase) GetUploadURL(ctx context.Context, userID uint32, crc32 string, sha1 string) (url string, err error) {
	return "", errors.New("todo")
}

func (uc *UserUseCase) BanUser(ctx context.Context, userID uint32) error {
	return uc.repo.BanUser(ctx, userID)
}

func (uc *UserUseCase) DeBanUser(ctx context.Context, userID uint32) error {
	return uc.repo.DeBanUser(ctx, userID)
}

func (uc *UserUseCase) EditUserTags(ctx context.Context, userID uint32, tags []uint32) (err error) {
	var info ModifiableUserInfo
	info = ModifiableUserInfo{
		Tags: util.NewValue[[]uint32](true, tags),
	}
	err = uc.repo.EditUserInfo(ctx, userID, &info)
	// TODO error
	return
}
