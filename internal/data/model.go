package data

import (
	"github.com/thgamejam/pkg/database"
)

// UserDB 用户模型
type UserDB struct {
	ID               uint32 `json:"id" gorm:"column:id"`
	AccountID        uint32 `json:"account_id" gorm:"column:account_id"`               // 账户id索引
	Name             string `json:"name" gorm:"column:name"`                           // 名称
	AvatarID         uint32 `json:"avatar_id" gorm:"column:avatar_id"`                 // 头像id
	Bio              string `json:"bio" gorm:"column:bio"`                             // 个人简介
	DisplayTag1      uint16 `json:"display_tag1" gorm:"column:display_tag1"`           // 展示的标签1
	DisplayTag2      uint16 `json:"display_tag2" gorm:"column:display_tag2"`           // 展示的标签2
	DisplayTag3      uint16 `json:"display_tag3" gorm:"column:display_tag3"`           // 展示的标签3
	AllowSyndication bool   `json:"allow_syndication" gorm:"column:allow_syndication"` // 是否允许联合发布邀请
	database.Model
}

func (UserDB) TableName() string {
	return "user"
}

// UserTagRelationalDB 用户与标签关系
type UserTagRelationalDB struct {
	ID     uint32 `json:"id" gorm:"column:id"`
	UserID uint32 `json:"user_id" gorm:"column:user_id"`         // 用户id
	TagID  uint16 `json:"user_tag_id" gorm:"column:user_tag_id"` // 用户标签索引
	Status uint16 `json:"status" gorm:"column:status"`           // 标签状态
	database.Model
}

func (UserTagRelationalDB) TableName() string {
	return "user_tag_relational"
}

// UserTagDB 用户标签模型
type UserTagDB struct {
	ID      uint16 `json:"id" gorm:"column:id"`
	Content string `json:"content" gorm:"column:content"` // 标签内容
	database.Model
}

func (UserTagDB) TableName() string {
	return "user_tag_enum"
}

// UserCache 用户信息缓存模型
type UserCache struct {
	ID        uint32   `json:"id"`         // 用户id
	AccountID uint32   `json:"account_id"` // 账户id
	Username  string   `json:"name"`       // 用户名
	Bio       string   `json:"bio"`        // 简介
	Tags      []string `json:"tags"`       // 标签
	AvatarUrl string   `json:"avatar_url"` // 头像链接
	AvatarExp int64    `json:"avatar_exp"` // 头像到期时间
}
