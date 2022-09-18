package data

import (
	"github.com/thgamejam/pkg/database"
)

// UserDB 用户模型
type UserDB struct {
	database.Model
	AccountID        uint32 `json:"account_id" gorm:"column:account_id; uniqueIndex"`  // 账户id索引
	Username         string `json:"uname" gorm:"column:uname"`                         // 名称
	AvatarID         uint32 `json:"avatar_id" gorm:"column:avatar_id"`                 // 头像id
	Bio              string `json:"bio" gorm:"column:bio"`                             // 个人简介
	DisplayTag1      uint16 `json:"display_tag1" gorm:"column:display_tag1"`           // 展示的标签1
	DisplayTag2      uint16 `json:"display_tag2" gorm:"column:display_tag2"`           // 展示的标签2
	DisplayTag3      uint16 `json:"display_tag3" gorm:"column:display_tag3"`           // 展示的标签3
	AllowSyndication bool   `json:"allow_syndication" gorm:"column:allow_syndication"` // 是否允许联合发布邀请
	// 用户状态, 二进制开关 UserStatus
	Status uint8 `json:"status" gorm:"column:status"`
}

func (UserDB) TableName() string {
	return "user"
}

// UserTagRelationalDB 用户与标签关系
type UserTagRelationalDB struct {
	database.Model
	database.DeleteModel
	UserID uint32 `json:"user_id" gorm:"column:user_id; index"`  // 用户id
	TagID  uint16 `json:"user_tag_id" gorm:"column:user_tag_id"` // 用户标签索引
}

func (UserTagRelationalDB) TableName() string {
	return "user_tag_relational"
}

// UserTagEnumDB 用户标签模型
type UserTagEnumDB struct {
	database.Model
	Content string `json:"content" gorm:"column:content"` // 标签内容
}

func (UserTagEnumDB) TableName() string {
	return "user_tag_enum"
}

// UserCache 用户信息缓存模型
type UserCache struct {
	ID        uint32   `json:"id"`         // 用户id
	AccountID uint32   `json:"account_id"` // 账户id
	Username  string   `json:"uname"`      // 用户名
	Bio       string   `json:"bio"`        // 简介
	Tags      []string `json:"tags"`       // 标签
	AvatarUrl string   `json:"image"`      // 头像
}

// TagLocalCache 标签本地缓存模型
type TagLocalCache struct {
	TagID    uint16 // 标签id
	Content  string // 标签内容
	IsExist  bool   // 是否存在
	QueryEXP int64  // 查询失效时间，时间戳
}

// Relationship 关系映射模型
type Relationship struct {
	UserID       uint32 `json:"user_id" gorm:"column:user_id; index"`
	FollowUserid uint32 `json:"followUserid" gorm:"column:follow_userid"`
	database.Model
}

func (Relationship) TableName() string {
	return "user_relationship"
}

// UserFollowInfo 用户关注信息模型
type UserFollowInfo struct {
	UserID      uint32 `json:"user_id" gorm:"column:user_id; index"`
	FansCount   uint32 `json:"fansCount" gorm:"column:fans_count"`
	FollowCount uint32 `json:"followCount" gorm:"column:follow_count"`
	database.Model
}

func (UserFollowInfo) TableName() string {
	return "user_follow_info"
}
