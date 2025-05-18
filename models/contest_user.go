package models

import (
	"gorm.io/gorm"
)

// ContestUser 表示竞赛用户的模型结构
// 该模型用于关联竞赛和用户，存储每个竞赛所包含的用户信息
type ContestUser struct {
	// ID 是该记录的主键，用于唯一标识每条竞赛用户关联记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// ContestId 表示竞赛的 ID，关联到竞赛表
	ContestId uint `gorm:"column:contest_id;type:int(11);" json:"contest_id"`
	// UserIdentity 表示用户的唯一标识，关联到用户表
	UserIdentity string `gorm:"column:user_identity;type:int(11);" json:"user_identity"`
	// UserBasic 是关联的用户基础信息，通过 user_identity 关联到 UserBasic 表
	UserBasic *UserBasic `gorm:"foreignKey:identity;references:user_identity;" json:"user_basic"`
}

// TableName 指定该模型对应的数据库表名
func (table *ContestUser) TableName() string {
	return "contest_user"
}
