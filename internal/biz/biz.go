package biz

import (
	"context"

	"github.com/google/wire"

	"github.com/thgamejam/pkg/util"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUseCase)

// UserRepo is a User repo.
type UserRepo interface {

	// GetUserStatus 获取用户状态
	GetUserStatus(ctx context.Context, userID []uint32) (map[uint32]*UserStatus, error)

	// GetUserInfoByAccountID 通过账户ID获取用户
	GetUserInfoByAccountID(ctx context.Context, accountID uint32) (user util.Val[*UserInfo], err error)

	// GetUserInfoByUserID 通过用户ID获取用户信息
	GetUserInfoByUserID(ctx context.Context, userID ...uint32) (user map[uint32]util.Val[*UserInfo], err error)

	// GetUserOwnTags 获取用户所有标签
	GetUserOwnTags(ctx context.Context, userID uint32) (tags []uint16, err error)

	// GetMultipleEnumTagContent 获取多个枚举列表中用户标签内容
	GetMultipleEnumTagContent(ctx context.Context, tagID ...uint16) ([]util.Val[*string], error)

	// GetUploadAvatarURL 获取头像上传链接
	GetUploadAvatarURL(ctx context.Context, userID uint32, crc32 string, sha1 string) (url string, err error)

	//CreateUser 创建用户
	CreateUser(ctx context.Context, accountID uint32, username string) (user *UserInfo, err error)

	// EditUserInfo 修改用户
	EditUserInfo(ctx context.Context, userID uint32, user *ModifiableUserInfo) error

	// EditUserStatus 修改用户状态
	EditUserStatus(ctx context.Context, userID uint32, userStatus *UserStatus) error

	// GetUserFansListByPage 获取用户粉丝分页列表
	GetUserFansListByPage(ctx context.Context, userid uint32, page uint32) (idList []uint32, err error)

	// GetUserFollowListByPage 获取用户关注分页列表
	GetUserFollowListByPage(ctx context.Context, userid uint32, page uint32) (isList []uint32, err error)

	// GetUserFollowInfo 获取用户关注信息
	GetUserFollowInfo(ctx context.Context, userid uint32) (followCount, fansCount uint32, err error)

	// AddUserFollowInfo 添加用户关注信息
	AddUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error)

	// DeleteUserFollowInfo 删除用户关注信息
	DeleteUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error)
}
