package models

import (
	"gorm.io/gorm"
)

// SubmitBasic 表示提交基础信息的模型结构
// 该模型用于存储用户提交代码的基本信息，以及关联的问题和用户信息
type SubmitBasic struct {
	// ID 是该记录的主键，用于唯一标识每条提交记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该提交记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该提交记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// Identity 是提交记录的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	// ProblemIdentity 表示该提交所属问题的唯一标识
	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	// ProblemBasic 是关联的问题基础信息，通过 problem_identity 关联到 ProblemBasic 表
	ProblemBasic *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity;" json:"problem_basic"`
	// UserIdentity 表示提交用户的唯一标识
	UserIdentity string `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	// UserBasic 是关联的用户基础信息，通过 user_identity 关联到 UserBasic 表
	UserBasic *UserBasic `gorm:"foreignKey:identity;references:user_identity;" json:"user_basic"`
	// Path 是提交代码的存放路径
	Path string `gorm:"column:path;type:varchar(255);" json:"path"`
	// Status 表示提交的状态，-1 表示待判断，1 表示答案正确，2 表示答案错误，3 表示运行超时，4 表示运行超内存，5 表示编译错误，6 表示非法代码
	Status int `gorm:"column:status;type:tinyint(1);" json:"status"`
}

// TableName 指定该模型对应的数据库表名
func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

// GetSubmitList 根据问题标识、用户标识和提交状态查询提交列表
// 返回一个 GORM 的查询构建器，可用于进一步的查询操作
func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	// 构建查询语句，预加载关联的问题和用户信息，并排除问题的内容和用户的密码
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	})
	// 如果问题标识不为空，添加问题标识的查询条件
	if problemIdentity != "" {
		tx.Where("problem_identity = ? ", problemIdentity)
	}
	// 如果用户标识不为空，添加用户标识的查询条件
	if userIdentity != "" {
		tx.Where("user_identity = ? ", userIdentity)
	}
	// 如果提交状态不为 0，添加提交状态的查询条件
	if status != 0 {
		tx.Where("status = ? ", status)
	}
	// 按提交记录的 ID 降序排序
	return tx.Order("submit_basic.id DESC")
}
