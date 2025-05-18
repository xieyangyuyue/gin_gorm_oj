package models

import (
	"gorm.io/gorm"
)

// ContestBasic 表示竞赛基础信息的模型结构
// 该模型用于存储竞赛的基本信息，以及关联的题目和用户信息
type ContestBasic struct {
	// ID 是该记录的主键，用于唯一标识每个竞赛记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该竞赛记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该竞赛记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// StartAt 是竞赛的开始时间
	StartAt MyTime `json:"start_at"`
	// EndAt 是竞赛的结束时间
	EndAt MyTime `json:"end_at"`
	// Identity 是竞赛的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	// Name 是竞赛的名称
	Name string `gorm:"column:name;type:varchar(100);" json:"name"`
	// Content 是竞赛的描述信息
	Content string `gorm:"column:content;type:text;" json:"content"`
	// ContestProblems 是关联的竞赛题目列表，通过 contest_id 关联到 ContestProblem 表
	ContestProblems []*ContestProblem `gorm:"foreignKey:contest_id;references:id;" json:"contest_problems"`
	// ContestUsers 是关联的竞赛用户列表，通过 contest_id 关联到 ContestUser 表
	ContestUsers []*ContestUser `gorm:"foreignKey:contest_id;references:id;" json:"contest_users"`
}

// TableName 指定该模型对应的数据库表名
func (table *ContestBasic) TableName() string {
	return "contest_basic"
}

// GetContestList 根据关键字查询竞赛列表
// 返回一个 GORM 的查询构建器，可用于进一步的查询操作
func GetContestList(keyword string) *gorm.DB {
	// 构建查询语句，选择竞赛的基本信息，并预加载关联的题目和题目基础信息
	tx := DB.Model(new(ContestBasic)).Distinct("`contest_basic`.`id`").Select("DISTINCT(`contest_basic`.`id`), `contest_basic`.`identity`, "+
		"`contest_basic`.`name`, `contest_basic`.`content`, `contest_basic`.`start_at`, `contest_basic`.`end_at`, `contest_basic`.`created_at`, `contest_basic`.`updated_at`, `contest_basic`.`deleted_at` ").Preload("ContestProblems").Preload("ContestProblems.ProblemBasic").
		Where("name like ? OR content like ? ", "%"+keyword+"%", "%"+keyword+"%")
	// 按竞赛 ID 降序排序
	return tx.Order("contest_basic.id DESC")
}
