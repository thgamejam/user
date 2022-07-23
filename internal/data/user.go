package data

import (
	"context"

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
func (r *userRepo) GetUserInfoByUserID(ctx context.Context, userID uint32) (user util.Val[*biz.UserInfo], err error) {

	user, err = r.cacheGetUser(ctx, userID)
	if err != nil {
		return
	}
	if user.IsExist() {
		return
	}

	// 如果缓存中不存在，再从数据库中获取
	user, err = r.dbGetUserByUserID(ctx, userID)
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
	err = r.cacheDelUserStatus(ctx, user.ID)
	if err != nil {
		// TODO log
	}
	return user, nil
}

// dbGetUserByUserID 在数据库中通过用户ID获取用户
func (r *userRepo) dbGetUserByUserID(
	ctx context.Context, userID uint32) (user util.Val[*biz.UserInfo], err error) {

	var model UserDB
	tx := r.data.sql.Limit(1).Find(&model, userID)
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
	ctx context.Context, accountID uint32) (user util.Val[*biz.UserInfo], err error) {

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

	return r.cacheGetUser(ctx, userID)
}

// cacheGetUser 在缓存中通过用户id获取用户信息
func (r *userRepo) cacheGetUser(ctx context.Context, id uint32) (user util.Val[*biz.UserInfo], err error) {
	var cache UserCache
	ok, err := r.data.rdb.Get(ctx, userCacheKey(id), &cache)
	if err != nil {
		return
	}
	if !ok {
		return
	}

	return util.NewValue(true, &biz.UserInfo{
		ID:        cache.ID,
		AccountID: cache.AccountID,
		Username:  cache.Username,
		Bio:       cache.Bio,
		Tags:      cache.Tags,
		AvatarUrl: cache.AvatarUrl,
	}), nil
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
