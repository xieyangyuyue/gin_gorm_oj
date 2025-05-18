package models

import (
	"gorm.io/gorm"
)

// ProblemCategory 表示问题分类的模型结构
// 该模型用于关联问题和分类，存储每个问题所属的分类信息
type ProblemCategory struct {
	// ID 是该记录的主键，用于唯一标识每条问题分类关联记录
	ID            uint           `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该记录的创建时间
	CreatedAt     MyTime         `json:"created_at"`
	// UpdatedAt 记录该记录的最后更新时间
	UpdatedAt     MyTime         `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt     gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// ProblemId 表示问题的 ID，关联到问题表
	ProblemId     uint           `gorm:"column:problem_id;type:int(11);" json:"problem_id"`
	// CategoryId 表示分类的 ID，关联到分类表
	CategoryId    uint           `gorm:"column:category_id;type:int(11);" json:"category_id"`
	// CategoryBasic 是关联的分类基础信息，通过 category_id 关联到 CategoryBasic 表
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id;" json:"category_basic"`
}

// TableName 指定该模型对应的数据库表名
func (table *ProblemCategory) TableName() string {
	return "problem_category"
}