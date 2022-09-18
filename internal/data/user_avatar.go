package data

import (
	"context"
	"time"

	"github.com/thgamejam/pkg/authentication"
	"github.com/thgamejam/pkg/util/strconv"
	"github.com/thgamejam/pkg/uuid"
)

const fileSuffix = ".webp"

// GetUploadAvatarURL 获取头像上传链接
func (r *userRepo) GetUploadAvatarURL(ctx context.Context, userID uint32, crc32 string, sha1 string) (url string, err error) {
	var claims authentication.UploadFileClaims
	claims = authentication.UploadFileClaims{
		Bucket:    r.conf.UserAvatarBucketName,
		Name:      strconv.UItoa[uint32](userID) + fileSuffix,
		UUID:      uuid.New(),
		ExpiresAt: uint32(time.Now().Unix()) + 3600, // TODO 头像到期时间要大于用户信息缓存到期时间。然而用户信息缓存到期时间没写啊
		CRC:       crc32,
		SHA1:      sha1,
	}

	uploadURL, err := authentication.CreateUploadURL(&claims, &r.conf.UploadFileServiceSecretKey)
	if err != nil {
		return "", err
	}
	return uploadURL, nil
}

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
	return r.conf.DefaultUserAvatarUrl
}
