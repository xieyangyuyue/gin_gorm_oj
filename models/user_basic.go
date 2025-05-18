package models

import (
	"gorm.io/gorm"
)

// UserBasic 表示用户基础信息的模型结构
// 该模型用于存储用户的基本信息
type UserBasic struct {
	// ID 是该记录的主键，用于唯一标识每个用户记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该用户记录的创建时间
	CreatedAt MyTime `gorm:"type:timestamp;" json:"created_at"`
	// UpdatedAt 记录该用户记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// Identity 是用户的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	// Name 是用户的姓名
	Name string `gorm:"column:name;type:varchar(100);" json:"name"`
	// Password 是用户的密码，不参与 JSON 序列化
	Password string `gorm:"column:password;type:varchar(32);" json:"-"`
	// Phone 是用户的手机号码
	Phone string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	// Mail 是用户的邮箱地址
	Mail string `gorm:"column:mail;type:varchar(100);" json:"mail"`
	// PassNum 是用户通过的问题数量
	PassNum int64 `gorm:"column:pass_num;type:int(11);" json:"pass_num"`
	// SubmitNum 是用户提交的问题数量
	SubmitNum int64 `gorm:"column:submit_num;type:int(11);" json:"submit_num"`
	// IsAdmin 表示用户是否为管理员，0 表示否，1 表示是
	IsAdmin int `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`
}

// TableName 指定该模型对应的数据库表名
func (table *UserBasic) TableName() string {
	return "user_basic"
}
