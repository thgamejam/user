package data

import (
	"context"
	v1 "user/proto/api/user/v1"
)

// GetUserTagContent 获取用户标签内容
func (r *userRepo) GetUserTagContent(ctx context.Context, tagID uint16) (*string, error) {
	// 补充长度
	if int(tagID) >= len(r.tagCache) {
		empty := make([]string, int(tagID)-len(r.tagCache)+1)
		r.tagCache = append(r.tagCache, empty...)
	}

	if r.tagCache[tagID] != "" {
		return &r.tagCache[tagID], nil
	}

	var tag UserTagDB
	ok, err := r.DBGetUserTag(ctx, &tag, tagID)
	if err != nil {
		return nil, err
	}
	if !ok {
		// 未找到标签
		return nil, v1.ErrorUserTag("not found user tag")
	}

	r.tagCache[tagID] = tag.Content
	return &r.tagCache[tagID], nil
}

// CreateUserTag 创建用户标签
func (r *userRepo) CreateUserTag(ctx context.Context, tagContent string) error {
	for _, content := range r.tagCache {
		if tagContent == content {
			// 无法创建重复的标签
			return v1.ErrorUserTag("duplicate user tag")
		}
	}

	var tag UserTagDB
	// 使用Limit(1)避免ErrRecordNotFound
	tx := r.data.DataBase.Where("content = ?", tagContent).Limit(1).Find(&tag)
	if tx.Error != nil {
		r.log.Errorf("CreateUserTag - DB.Where().First - err=%v", tx.Error)
		return v1.ErrorInternalServer("database error")
	}
	if tx.RowsAffected != 0 {
		return v1.ErrorUserTag("duplicate user tag")
	}

	// 创建标签
	tag = UserTagDB{
		Content: tagContent,
	}
	err := r.data.DataBase.Create(&tag).Error
	if err != nil {
		r.log.Errorf("CreateUserTag - DB.Create - err=%v", err)
		return v1.ErrorInternalServer("database error")
	}

	return nil
}

// DBGetUserTag 从数据库中获取标签
func (r *userRepo) DBGetUserTag(ctx context.Context, model *UserTagDB, id uint16) (bool, error) {
	tx := r.data.DataBase.Limit(1).Find(model, id)
	if tx.Error != nil {
		r.log.Errorf("DBGetUserTag - DB.Limit(1).Find - err=%v", tx.Error)
		return false, v1.ErrorInternalServer("database error")
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// DBGetALLUserTagContent 缓存所有的标签内容
func (r *userRepo) DBGetALLUserTagContent(ctx context.Context) ([]string, error) {
	var tags []UserTagDB
	tags = []UserTagDB{}

	tx := r.data.DataBase.Order("id").Find(&tags)
	if tx.Error != nil {
		r.log.Errorf("DBGetALLUserTagContent - DB.Order(id).Find - err=%v", tx.Error)
		return nil, v1.ErrorInternalServer("database error")
	}

	length := tags[len(tags)-1].ID
	tagStr := make([]string, length, length)

	for _, tag := range tags {
		tagStr[tag.ID] = tag.Content
	}

	return tagStr, nil
}
