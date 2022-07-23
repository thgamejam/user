package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/thgamejam/pkg/cache"
	"github.com/thgamejam/pkg/database"
	"github.com/thgamejam/pkg/object_storage"
	"gorm.io/gorm"

	"user/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
)

// Data .
type Data struct {
	rdb *cache.Cache
	sql *gorm.DB
	oss *object_storage.ObjectStorage
}

// NewData .
func NewData(confData *conf.Data, confUser *conf.User, logger log.Logger) (*Data, func(), error) {
	redis, err := cache.NewCache(confData.Redis)
	if err != nil {
		return nil, nil, err
	}
	db, err := database.NewDataBase(confData.Database)
	if err != nil {
		return nil, nil, err
	}
	oss, err := object_storage.NewObjectStorage(confData.ObjectStorage)
	if err != nil {
		return nil, nil, err
	}
	data := &Data{
		rdb: redis,
		sql: db,
		oss: oss,
	}

	ctx := context.Background()
	ok, err := data.oss.ExistBucket(ctx, confUser.UserAvatarBucketName)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		err := data.oss.CreateBucket(ctx, confUser.DefaultUserAvatarKey)
		if err != nil {
			return nil, nil, err
		}
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return data, cleanup, nil
}
