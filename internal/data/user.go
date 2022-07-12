package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/thgamejam/pkg/util/strconv"
	"user/internal/conf"

	"user/internal/biz"
)

type userRepo struct {
	data     *Data
	conf     *conf.User
	tagCache []string // 用户标签缓存

	log *log.Helper
}

var userCacheKey = func(id uint32) string {
	return "user_model_" + strconv.UItoa(id)
}

// 缓存索引 account_id => user_id
var userIndexByAccountIDCacheKey = func(id uint32) string {
	return "user_idx_account_" + strconv.UItoa(id)
}

var userAvatarIDCacheURL = func(hash string) string {
	return "user_AvatarID_URL_" + hash
}

// NewUserRepo .
func NewUserRepo(conf *conf.User, data *Data, logger log.Logger) (biz.UserRepo, error) {
	r := &userRepo{
		data: data,
		conf: conf,
		log:  log.NewHelper(logger),
	}
	ctx := context.Background()
	// 读取所有的用户标签，缓存到本地
	tags, err := r.DBGetALLUserTagContent(ctx)
	if err != nil {
		return nil, err
	}
	r.tagCache = tags
	return r, nil
}
