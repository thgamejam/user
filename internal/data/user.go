package data

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/thgamejam/pkg/util"
	"github.com/thgamejam/pkg/util/strconv"
	"user/internal/biz"
)

var userCacheKey = func(userID uint32) string {
	return "user_model_" + strconv.UItoa(userID)
}

// 缓存索引 account_id => user_id
var userIndexByAccountIDCacheKey = func(accountID uint32) string {
	return "user_idx_account_" + strconv.UItoa(accountID)
}

// GetUserInfoByUserID 通过用户ID获取用户信息
func (r *userRepo) GetUserInfoByUserID(
	ctx context.Context, userID ...uint32) (res map[uint32]util.Val[*biz.UserInfo], err error) {

	res, err = r.cacheGetUser(ctx, userID...)
	if err != nil {
		return
	}

	// 如果缓存中不存在，再从数据库中获取
	undone := make([]uint32, len(userID))
	for _, id := range userID {
		if info, ok := res[id]; ok {
			if !info.IsExist() {
				undone = append(undone, id)
			}
		} else {
			undone = append(undone, id)
		}
	}

	models, err := r.dbGetUserByUserID(ctx, userID...)
	if err != nil {
		return
	}

	status := make(map[uint32]uint8)
	for id, v := range models {
		if !v.IsExist() {
			// 数据库中不存在
			status[id] = 0
			continue
		}

		model := v.Val()
		// 获取标签内容
		tagList, err := r.GetMultipleEnumTagContent(ctx, model.DisplayTag1, model.DisplayTag2, model.DisplayTag3)
		if err != nil {
			return
		}
		// 清除不存在的标签id
		tags := make([]string, 0, 3)
		for _, v := range tagList {
			if v.IsExist() {
				tags = append(tags, *v.Val())
			}
		}

		avatarURL, err := r.GetUserAvatarURL(ctx, model.AvatarID)
		if err != nil {
			// TODO 这里不应该发生错误，但还是需要处理
		}

		var info biz.UserInfo
		info = biz.UserInfo{
			ID:        model.ID,
			AccountID: model.AccountID,
			Username:  model.Username,
			Bio:       model.Bio,
			Tags:      tags,
			AvatarUrl: avatarURL,
		}
		res[model.ID] = util.NewValue(true, &info)

		// 保存到缓存中
		err = r.cacheSetUser(ctx, &info)
		if err != nil {
			// TODO 应该打印错误
		}
		status[id] = model.Status
	}

	// 用户状态保存到缓存
	err = r.cacheSetUserStatus(ctx, status)
	if err != nil {
		// TODO 应该打印错误
	}

	return
}

// GetUserInfoByAccountID 通过账户ID获取用户信息
func (r *userRepo) GetUserInfoByAccountID(
	ctx context.Context, accountID uint32) (user util.Val[*biz.UserInfo], err error) {

	// 从缓存中获取
	user, err = r.cacheGetUserByAccountID(ctx, accountID)
	if err != nil {
		return
	}
	// 缓存中存在，直接返回
	if user.IsExist() {
		return
	}

	// 如果缓存中不存在，再从数据库中获取
	user, err = r.dbGetUserByAccountID(ctx, accountID)
	if err != nil {
		return
	}
	// 数据库中不存在
	if !user.IsExist() {
		return
	}

	// 保存到缓存中
	err = r.cacheSetUser(ctx, user.Val())
	if err != nil {
		// TODO log
	}

	return
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, accountID uint32, username string) (user *biz.UserInfo, err error) {
	user, err = r.dbCreateUser(ctx, accountID, username)
	if err != nil {
		return
	}

	// 删除用户状态缓存
	err = r.cacheDelUserStatus(ctx, user.ID)
	if err != nil {
		// TODO log
	}
	return user, nil
}

func (r *userRepo) EditUserInfo(ctx context.Context, userID uint32, info *biz.ModifiableUserInfo) (err error) {
	update := make(map[string]interface{}, 3)
	if info.Username.IsExist() {
		update["uname"] = info.Username.Val()
	}
	if info.Bio.IsExist() {
		update["bio"] = info.Bio.Val()
	}
	if info.Tags.IsExist() {
		update["display_tag1"] = ""
		update["display_tag2"] = ""
		update["display_tag3"] = ""
		size := len(info.Tags.Val())
		if size > 0 {
			update["display_tag1"] = info.Tags.Val()[0]
		} else if size > 1 {
			update["display_tag2"] = info.Tags.Val()[1]
		} else if size > 2 {
			update["display_tag3"] = info.Tags.Val()[2]
		}
	}
	tx := r.data.sql.Model(&UserDB{}).Where("id = ?", userID).Updates(update)
	if err = tx.Error; err != nil {
		return
	}

	err = r.cacheDelUser(ctx, userID)
	return
}

// dbGetUserByUserID 在数据库中通过用户ID获取用户
func (r *userRepo) dbGetUserByUserID(
	ctx context.Context, userID ...uint32) (res map[uint32]util.Val[*UserDB], err error) {

	// 初始化
	notExist := util.NewValue[*UserDB](false, nil)
	length := len(userID)
	res = make(map[uint32]util.Val[*UserDB], length)
	for _, id := range userID {
		res[id] = notExist
	}

	models := make([]UserDB, length)
	tx := r.data.sql.Find(&models, userID)
	if tx.Error != nil {
		err = tx.Error
		// TODO 处理错误
	}
	if tx.RowsAffected == 0 {
		return
	}

	// 填充数据
	for _, model := range models {
		res[model.ID] = util.NewValue(true, &model)
	}

	return res, nil
}

// dbGetUserByAccountID 在数据库中通过账户ID获取用户
func (r *userRepo) dbGetUserByAccountID(
	ctx context.Context, accountID uint32) (user util.Val[*biz.UserInfo], err error) {

	var model UserDB
	tx := r.data.sql.Limit(1).Find(&model, "account_id = ?", accountID)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	if tx.RowsAffected == 0 {
		return
	}

	// 获取标签内容
	tagList, err := r.GetMultipleEnumTagContent(ctx, model.DisplayTag1, model.DisplayTag2, model.DisplayTag3)
	if err != nil {
		return
	}
	// 清除不存在的标签id
	tags := make([]string, 0, 3)
	for _, v := range tagList {
		if v.IsExist() {
			tags = append(tags, *v.Val())
		}
	}

	avatarURL, err := r.GetUserAvatarURL(ctx, model.AvatarID)
	if err != nil {
		return
	}

	return util.NewValue(true, &biz.UserInfo{
		ID:        model.ID,
		AccountID: model.AccountID,
		Username:  model.Username,
		Bio:       model.Bio,
		Tags:      tags,
		AvatarUrl: avatarURL,
	}), nil
}

// dbCreateUser 在数据库中创建用户信息
func (r *userRepo) dbCreateUser(
	ctx context.Context, accountID uint32, username string) (user *biz.UserInfo, err error) {

	model := UserDB{
		AccountID: accountID,
		Username:  username,
	}

	err = r.data.sql.Create(&model).Error
	if err != nil {
		return
	}

	return &biz.UserInfo{
		ID:        model.ID,
		AccountID: model.AccountID,
		Username:  model.Username,
		Bio:       model.Bio,
		Tags:      []string{},
		AvatarUrl: r.GetDefaultUserAvatarURL(),
	}, nil
}

// cacheGetUserByAccountID 在缓存中通过账户ID获取用户信息
func (r *userRepo) cacheGetUserByAccountID(
	ctx context.Context, accountID uint32) (info util.Val[*biz.UserInfo], err error) {

	// 从缓存中使用account-id获取user-id
	userStrID, ok, err := r.data.rdb.GetString(ctx, userIndexByAccountIDCacheKey(accountID))
	if err != nil {
		return
	}
	if !ok {
		return
	}

	var userID uint32
	err = strconv.UParse(userStrID, &userID)
	if err != nil {
		return
	}

	res, err := r.cacheGetUser(ctx, userID)
	if err != nil {
		return
	}
	info, ok = res[userID]
	if !ok {
		// TODO 没有对应用户id的数据，这不应该发生，这样应该触发错误
	}
	return info, nil
}

// cacheGetUser 在缓存中通过用户id获取用户信息
func (r *userRepo) cacheGetUser(
	ctx context.Context, userID ...uint32) (user map[uint32]util.Val[*biz.UserInfo], err error) {

	length := len(userID)
	pipe := r.data.rdb.Client.Pipeline()
	cmds := make(map[uint32]*redis.StringCmd, length)
	for _, id := range userID {
		cmds[id] = pipe.Get(ctx, userCacheKey(id))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	caches := make([]util.Val[UserCache], length, length)
	for id, cmd := range cmds {
		// 判断返回是空
		if cmd.Err() == redis.Nil {
			caches[id] = util.NewValue[UserCache](false, UserCache{})
			continue
		}

		// 这里不应该发生错误，应该在Exec时返回错误
		if cmd.Err() != nil {
			return nil, cmd.Err()
		}

		v, err := cmd.Result()
		if err != nil {
			return nil, err
		}

		var info UserCache
		info = UserCache{}
		err = json.Unmarshal([]byte(v), &info)
		if err != nil {
			return nil, err
		}
		caches[id] = util.NewValue[UserCache](false, info)
	}

	res := make(map[uint32]util.Val[*biz.UserInfo], length)
	for i, v := range caches {
		if v.IsExist() {
			cache := v.Val()
			res[cache.ID] = util.NewValue(true, &biz.UserInfo{
				ID:        cache.ID,
				AccountID: cache.AccountID,
				Username:  cache.Username,
				Bio:       cache.Bio,
				Tags:      cache.Tags,
				AvatarUrl: cache.AvatarUrl,
			})
		} else {
			res[userID[i]] = util.NewValue[*biz.UserInfo](false, nil)
		}
	}
	return res, nil
}

// cacheSetUser 在缓存中保存用户信息
func (r *userRepo) cacheSetUser(ctx context.Context, user *biz.UserInfo) error {

	err := r.data.rdb.SetString(ctx, userIndexByAccountIDCacheKey(user.AccountID), strconv.UItoa(user.ID), 0)
	if err != nil {
		return err
	}

	cache := UserCache{
		ID:        user.ID,
		AccountID: user.AccountID,
		Username:  user.Username,
		Bio:       user.Bio,
		Tags:      user.Tags,
		AvatarUrl: user.AvatarUrl,
	}
	err = r.data.rdb.Set(ctx, userCacheKey(user.ID), &cache, 0)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) cacheDelUser(ctx context.Context, userID uint32) error {
	return r.data.rdb.Del(ctx, userCacheKey(userID))
}
