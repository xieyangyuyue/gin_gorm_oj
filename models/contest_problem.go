package models

import (
	"gorm.io/gorm"
)

// ContestProblem 表示竞赛问题的模型结构
// 该模型用于关联竞赛和问题，存储每个竞赛所包含的问题信息
type ContestProblem struct {
	// ID 是该记录的主键，用于唯一标识每条竞赛问题关联记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// ContestId 表示竞赛的 ID，关联到竞赛表
	ContestId uint `gorm:"column:contest_id;type:int(11);" json:"contest_id"`
	// ProblemId 表示问题的 ID，关联到问题表
	ProblemId uint `gorm:"column:problem_id;type:int(11);" json:"problem_id"`
	// ProblemBasic 是关联的问题基础信息，通过 problem_id 关联到 ProblemBasic 表
	ProblemBasic *ProblemBasic `gorm:"foreignKey:id;references:problem_id;" json:"problem_basic"`
}

// TableName 指定该模型对应的数据库表名
func (table *ContestProblem) TableName() string {
	return "contest_problem"
}
