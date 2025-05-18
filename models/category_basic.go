package models

import (
	"gorm.io/gorm"
)

// CategoryBasic 表示分类基础信息的模型结构
// 该模型用于存储问题分类的基本信息
type CategoryBasic struct {
	// ID 是该记录的主键，用于唯一标识每个分类记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该分类记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该分类记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// Identity 是分类的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	// Name 是分类的名称
	Name string `gorm:"column:name;type:varchar(100);" json:"name"`
	// ParentId 是该分类的父级分类 ID
	ParentId int `gorm:"column:parent_id;type:int(11);" json:"parent_id"`
}

// TableName 指定该模型对应的数据库表名
func (table *CategoryBasic) TableName() string {
	return "category_basic"
}
