package data

import (
	"context"
	"github.com/thgamejam/pkg/util/strconv"
	"gorm.io/gorm"
)

var userFollowCacheKey = func(userid uint32) string {
	return "userFollowInfo_by_id" + strconv.Itoa(int(userid))
}

func (r *userRepo) GetUserFansListByPage(ctx context.Context, userid uint32, page uint32) (idList []uint32, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *userRepo) GetUserFollowListByPage(ctx context.Context, userid uint32, page uint32) (isList []uint32, err error) {
	//TODO implement me
	panic("implement me")
}

// AddUserFollowInfo 增加用户关系
func (r *userRepo) AddUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error) {
	// 查找数据库中用户是否已经关注
	yes, err := r.FindRelationship(ctx, userid, followUserId)
	if err != nil {
		return err
	}
	if yes {
		return nil
	}

	// 查找用户信息是否存在
	followC, fansC, err := r.GetUserFollowInfo(ctx, userid)
	if err != nil {
		return err
	}

	model := Relationship{
		UserID:       userid,
		FollowUserid: followUserId,
	}
	// 开启事物增加信息
	err = r.data.sql.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(model).Error; err != nil {
			return err
		}
		if err := tx.Updates(map[string]interface{}{"follow_count": followC + 1}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = r.data.rdb.Set(ctx, userFollowCacheKey(userid), &UserFollowInfo{
		UserID:      userid,
		FansCount:   fansC,
		FollowCount: followC + 1,
	}, 0)
	return nil
}

// DeleteUserFollowInfo 删除用户关系
func (r *userRepo) DeleteUserFollowInfo(ctx context.Context, userid uint32, followUserId uint32) (err error) {
	// 查找数据库中用户是否已经关注
	yes, err := r.FindRelationship(ctx, userid, followUserId)
	if err != nil {
		return err
	}
	if !yes {
		return nil
	}

	// 查找用户信息是否存在
	followC, fansC, err := r.GetUserFollowInfo(ctx, userid)
	if err != nil {
		return err
	}

	model := Relationship{
		UserID:       userid,
		FollowUserid: followUserId,
	}
	// 开启事物增加信息
	err = r.data.sql.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(model).Error; err != nil {
			return err
		}
		if err := tx.Updates(map[string]interface{}{"fans_count": fansC - 1}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	_ = r.data.rdb.Set(ctx, userFollowCacheKey(userid), &UserFollowInfo{
		UserID:      userid,
		FansCount:   fansC - 1,
		FollowCount: followC,
	}, 0)
	return nil
}

// GetUserFollowInfo 获取用户关注信息
func (r *userRepo) GetUserFollowInfo(ctx context.Context, userid uint32) (followCount, fansCount uint32, err error) {
	// 在缓存中查找
	followInfo := &UserFollowInfo{}
	ok, err := r.data.rdb.Get(ctx, userFollowCacheKey(userid), followInfo)
	if err != nil {
		return 0, 0, err
	}
	if ok {
		return followInfo.FollowCount, followInfo.FansCount, nil
	}

	// 若不存在则在数据库中查找
	var model UserFollowInfo
	tx := r.data.sql.Limit(1).Find(&model, "userid = ?", userid)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	// 数据不存在则创建用户数据信息
	if tx.RowsAffected == 0 {
		r.data.sql.Create(UserFollowInfo{
			UserID:      userid,
			FansCount:   0,
			FollowCount: 0,
		})
		_ = r.data.rdb.Set(ctx, userFollowCacheKey(userid), &model, 0)
		return 0, 0, nil
	}
	_ = r.data.rdb.Set(ctx, userFollowCacheKey(userid), &model, 0)

	return model.FollowCount, model.FansCount, nil
}

// FindRelationship 获取用户是否关注另一个用户
func (r *userRepo) FindRelationship(ctx context.Context, userid, followUserId uint32) (exist bool, err error) {
	var model Relationship
	tx := r.data.sql.Limit(1).Find(&model, "userid = ?", userid)
	if tx.Error != nil {
		err = tx.Error
		return false, err
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
