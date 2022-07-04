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
	Cache         *cache.Cache
	DataBase      *gorm.DB
	ObjectStorage *object_storage.ObjectStorage
	UserConf      *conf.User
}

// NewData .
func NewData(c *conf.Data, conf *conf.User, logger log.Logger) (*Data, func(), error) {
	newCache, _ := cache.NewCache(c.Redis)
	newDataBase, _ := database.NewDataBase(c.Database)
	newObjectStorage, _ := object_storage.NewObjectStorage(c.ObjectStorage)
	data := &Data{
		Cache:         newCache,
		DataBase:      newDataBase,
		ObjectStorage: newObjectStorage,
		UserConf:      conf,
	}

	ctx := context.Background()
	ok, err := data.ObjectStorage.ExistBucket(ctx, conf.UserAvatarBucketName)
	if err != nil {
		return nil, nil, err
	}
	if err == nil {
		if !ok {
			err := data.ObjectStorage.CreateBucket(ctx, conf.DefaultUserAvatarHash)
			if err != nil {
				//panic(err)
			}
		}
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return data, cleanup, nil
}
