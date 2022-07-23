package data

import (
	"context"
	"github.com/thgamejam/pkg/util/strconv"
)

const fileSuffix = ".webp"

var _defaultUserAvatarURL string

// GetUserAvatarURL 根据用户ID获取头像
func (r *userRepo) GetUserAvatarURL(ctx context.Context, avatarID uint32) (string, error) {
	strID := strconv.UItoa(avatarID)

	// 从对象存储里构建新的URL
	fileName := strID + fileSuffix
	url, err := r.data.oss.PreSignGetURL(ctx, r.conf.UserAvatarBucketName, fileName, fileName, -1)
	if err != nil {
		r.log.Errorf("GetUserAvatarURL - OSS.PreSignGetURL(fileName) - err=%v", err)
		// 若对象存储报错则尝试获取默认头像
		return r.GetDefaultUserAvatarURL(), nil
	}

	return url.String(), nil
}

// GetDefaultUserAvatarURL 获取默认用户头像连接
func (r *userRepo) GetDefaultUserAvatarURL() string {
	return _defaultUserAvatarURL
}

// renewDefaultUserAvatarURL 重置默认用户头像url
func (r *userRepo) renewDefaultUserAvatarURL(ctx context.Context) (err error) {
	url, err := r.data.oss.PreSignGetURL(ctx,
		r.conf.UserAvatarBucketName,
		r.conf.DefaultUserAvatarKey,
		r.conf.DefaultUserAvatarKey,
		-1,
	)
	if err != nil {
		r.log.Errorf("GetUserAvatarURL - OSS.PreSignGetURL(r.conf.DefaultUserAvatarKey) - err=%v", err)
		return
	}
	_defaultUserAvatarURL = url.String()
	return
}
