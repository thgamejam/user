package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"user/internal/biz"
	"user/internal/conf"
)

type userRepo struct {
	data     *Data
	conf     *conf.User
	tagCache []*TagLocalCache // 用户标签缓存

	log *log.Helper
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
	err := r.localCacheSyncTags(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}
