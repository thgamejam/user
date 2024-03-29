package data

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/thgamejam/pkg/util"
	"github.com/thgamejam/pkg/util/strconv"
	"user/internal/biz"
)

type dataStatus util.Val[uint8]

func newDataStatus(inDB bool, value uint8) dataStatus {
	return util.NewValue(inDB, value)
}

var userStatusCacheKey = func(id uint32) string {
	return "user_status_" + strconv.UItoa(id)
}

// GetUserStatus 获取用户状态
// 如果不返回错误，则返回的map中一定有全量的用户状态，没有nil的情况
func (r *userRepo) GetUserStatus(ctx context.Context, userID []uint32) (map[uint32]*biz.UserStatus, error) {
	result, err := r.cacheGetUserStatus(ctx, userID)
	if err != nil {
		return nil, err
	}

	for id, v := range result {
		if v.IsExist() {
			continue
		}

		// 如果在redis中没有获取到，则从数据库中获取
		s, err := r.dbGetUserStatus(ctx, id)
		if err != nil {
			return nil, err
		}
		if s.IsExist() {
			result[id] = newDataStatus(true, s.Val())
			// 将数据存入缓存内
			_ = r.cacheSetUserStatus(ctx, map[uint32]uint8{id: s.Val()})
		}
	}

	status := make(map[uint32]*biz.UserStatus, len(userID))
	for _, id := range userID {
		s, ok := result[id]
		if !ok {
			// TODO 不应该发生的
		}
		status[id] = biz.NewUserStatus(s.Val())
	}

	return status, nil
}

func (r *userRepo) EditUserStatus(ctx context.Context, userID uint32, userStatus *biz.UserStatus) (err error) {
	err = r.dbSetUserStatus(ctx, userID, uint8(*userStatus))
	if err != nil {
		return
	}

	err = r.cacheDelUserStatus(ctx, userID)
	return
}

// dbGetUserStatus 从数据库中获取用户状态
func (r *userRepo) dbGetUserStatus(ctx context.Context, userID uint32) (status dataStatus, err error) {
	var user UserDB
	tx := r.data.sql.Limit(1).Select("status").Find(&user, userID)
	if tx.Error != nil {
		return status, tx.Error
	}
	if tx.RowsAffected == 0 {
		return
	}

	status = newDataStatus(true, user.Status)
	return
}

func (r *userRepo) dbSetUserStatus(ctx context.Context, userID uint32, status uint8) (err error) {
	tx := r.data.sql.Model(&UserDB{}).Where("id = ?", userID).Update("status", status)
	err = tx.Error
	return
}

func (r *userRepo) cacheGetUserStatus(ctx context.Context, userID []uint32) (status map[uint32]dataStatus, err error) {
	// 初始化数据
	length := len(userID)
	status = make(map[uint32]dataStatus, length)

	pipe := r.data.rdb.Client.Pipeline()
	cmds := make(map[uint32]*redis.StringCmd, length)
	for _, id := range userID {
		cmds[id] = pipe.Get(ctx, userStatusCacheKey(id))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	for id, cmd := range cmds {
		// 判断返回是空
		if cmd.Err() == redis.Nil {
			status[id] = newDataStatus(false, 0)
			continue
		}

		// 这里不应该发生错误，应该在Exec时返回错误
		if cmd.Err() != nil {
			return nil, cmd.Err()
		}

		v, err := cmd.Uint64()
		if err != nil {
			return nil, err
		}

		status[id] = newDataStatus(true, uint8(v))
	}

	return
}

func (r *userRepo) cacheSetUserStatus(ctx context.Context, status map[uint32]uint8) error {
	pipe := r.data.rdb.Client.Pipeline()
	for id, v := range status {
		pipe.Set(ctx, userStatusCacheKey(id), strconv.UItoa(v), 0)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) cacheDelUserStatus(ctx context.Context, userID uint32) error {
	err := r.data.rdb.Del(ctx, userStatusCacheKey(userID))
	if err != nil {
		return err
	}
	return nil
}
