package data

import (
	"context"
	"time"

	"github.com/thgamejam/pkg/util"
)

// GetMultipleEnumTagContent 获取多个枚举列表中用户标签内容
// 返回的字符串指针可能为nil。当返回字符串指针为nil时，表示tag-id对应的标签不存在
func (r *userRepo) GetMultipleEnumTagContent(ctx context.Context, tagID ...uint16) ([]util.Val[*string], error) {
	contents := make([]util.Val[*string], len(tagID), len(tagID))
	for i, id := range tagID {
		content, err := r.GetEnumTagContent(ctx, id)
		if err != nil {
			return nil, err
		}
		contents[i] = content
	}
	return contents, nil
}

// GetEnumTagContent 获取枚举列表中用户标签内容
func (r *userRepo) GetEnumTagContent(ctx context.Context, tagID uint16) (util.Val[*string], error) {
	// 如果是不会过期的缓存，则直接返回
	if r.tagCache[tagID].QueryEXP == -1 {
		return util.NewValue(r.tagCache[tagID].IsExist, &r.tagCache[tagID].Content), nil
	}

	// 如果缓存不需要重新获取，则直接返回
	t := time.Now().Unix()
	if t > r.tagCache[tagID].QueryEXP {
		return util.NewValue(r.tagCache[tagID].IsExist, &r.tagCache[tagID].Content), nil
	}

	// 从数据库中重新查询数据
	content, err := r.dbGetEnumTag(ctx, tagID)
	if err != nil {
		return nil, err
	}

	if !content.IsExist() {
		// 如果未找到标签，则记录到本地缓存，修改缓存查询到期时间
		r.tagCache[tagID].QueryEXP = t + 3600
		return nil, nil
	}

	// 找到标签，填充数据到本地缓存
	r.tagCache[tagID].TagID = tagID
	r.tagCache[tagID].Content = content.Val()
	r.tagCache[tagID].IsExist = true
	r.tagCache[tagID].QueryEXP = -1

	return util.NewValue(r.tagCache[tagID].IsExist, &r.tagCache[tagID].Content), nil
}

// CreateEnumTag 枚举列表中创建用户标签
func (r *userRepo) CreateEnumTag(ctx context.Context, tagContent string) error {
	var tag UserTagEnumDB
	// 创建标签
	tag = UserTagEnumDB{
		Content: tagContent,
	}
	err := r.data.sql.Create(&tag).Error
	if err != nil {
		return err
	}

	r.tagCache[tag.ID].TagID = uint16(tag.ID)
	r.tagCache[tag.ID].Content = tag.Content
	r.tagCache[tag.ID].IsExist = true
	r.tagCache[tag.ID].QueryEXP = -1

	return nil
}

// dbGetEnumTag 从数据库中获取标签
func (r *userRepo) dbGetEnumTag(ctx context.Context, id uint16) (content util.Val[string], err error) {
	var model UserTagEnumDB
	tx := r.data.sql.Limit(1).Find(&model, id)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	if tx.RowsAffected == 0 {
		return
	}
	return util.NewValue(true, model.Content), nil
}

// localCacheSyncTags 标签本地缓存同步数据
func (r *userRepo) localCacheSyncTags(ctx context.Context) error {
	// 初始化切片
	r.tagCache = make([]*TagLocalCache, 65535, 65535)
	t := time.Now().Unix() + 3600 // 使数据在初始化后一个小时再更新查询
	for i := uint16(0); i < 65535; i++ {
		r.tagCache[i] = &TagLocalCache{
			TagID:    i,
			Content:  "",
			IsExist:  false,
			QueryEXP: t,
		}
	}

	// 迭代器获取数据库数据
	rows, err := r.data.sql.Model(&UserTagEnumDB{}).Order("id").Rows()
	if err != nil {
		return err
	}
	var tag UserTagEnumDB
	for rows.Next() {
		err = r.data.sql.ScanRows(rows, &tag)
		if err != nil {
			return err
		}

		if "" != tag.Content {
			r.tagCache[tag.ID].Content = tag.Content
			r.tagCache[tag.ID].IsExist = true
		}
	}

	// 将已存在的id设置为不过期不重复查询
	for i := uint32(0); i <= tag.ID; i++ {
		r.tagCache[i].QueryEXP = -1
	}

	return nil
}
