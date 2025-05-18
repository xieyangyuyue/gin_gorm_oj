package models

import (
	"gorm.io/gorm"
)

// TestCase 表示测试用例的模型结构
// 该模型用于存储每个问题的测试用例信息
type TestCase struct {
	// ID 是该记录的主键，用于唯一标识每条测试用例记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// Identity 是测试用例的唯一标识
	Identity string `gorm:"column:identity;type:varchar(255);" json:"identity"`
	// ProblemIdentity 表示该测试用例所属问题的唯一标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(255);" json:"problem_identity"`
	// Input 是测试用例的输入内容
	Input string `gorm:"column:input;type:text;" json:"input"`
	// Output 是测试用例的预期输出内容
	Output string `gorm:"column:output;type:text;" json:"output"`
}

// TableName 指定该模型对应的数据库表名
func (table *TestCase) TableName() string {
	return "test_case"
}
