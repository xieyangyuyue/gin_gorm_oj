package models

import (
	"gorm.io/gorm"
)

// ProblemBasic 表示问题基础信息的模型结构
// 该模型用于存储问题的基本信息，以及关联的分类和测试用例信息
type ProblemBasic struct {
	// ID 是该记录的主键，用于唯一标识每个问题记录
	ID uint `gorm:"primarykey;" json:"id"`
	// CreatedAt 记录该问题记录的创建时间
	CreatedAt MyTime `json:"created_at"`
	// UpdatedAt 记录该问题记录的最后更新时间
	UpdatedAt MyTime `json:"updated_at"`
	// DeletedAt 是软删除标记，使用 gorm 的软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	// Identity 是问题的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	// ProblemCategories 是关联的问题分类列表，通过 problem_id 关联到 ProblemCategory 表
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"`
	// Title 是问题的标题
	Title string `gorm:"column:title;type:varchar(255);" json:"title"`
	// Content 是问题的详细内容
	Content string `gorm:"column:content;type:text;" json:"content"`
	// MaxRuntime 是问题的最大运行时长
	MaxRuntime int `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`
	// MaxMem 是问题的最大运行内存
	MaxMem int `gorm:"column:max_mem;type:int(11);" json:"max_mem"`
	// TestCases 是关联的测试用例列表，通过 problem_identity 关联到 TestCase 表
	TestCases []*TestCase `gorm:"foreignKey:problem_identity;references:identity;" json:"test_cases"`
	// PassNum 是问题的通过次数
	PassNum int64 `gorm:"column:pass_num;type:int(11);" json:"pass_num"`
	// SubmitNum 是问题的提交次数
	SubmitNum int64 `gorm:"column:submit_num;type:int(11);" json:"submit_num"`
}

// TableName 指定该模型对应的数据库表名
func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

// GetProblemList 根据关键字和分类标识查询问题列表
// 返回一个 GORM 的查询构建器，可用于进一步的查询操作
func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	// 构建查询语句，选择问题的基本信息，并预加载关联的分类和分类基础信息
	tx := DB.Model(new(ProblemBasic)).Distinct("`problem_basic`.`id`").
		Select("DISTINCT(`problem_basic`.`id`), `problem_basic`.`identity`, "+
			"`problem_basic`.`title`, `problem_basic`.`max_runtime`, `problem_basic`.`max_mem`, `problem_basic`.`pass_num`, "+
			"`problem_basic`.`submit_num`, `problem_basic`.`created_at`, `problem_basic`.`updated_at`, `problem_basic`.`deleted_at` ").
		Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ? ", "%"+keyword+"%", "%"+keyword+"%")
	// 如果分类标识不为空，添加分类标识的查询条件
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	// 按问题记录的 ID 降序排序
	return tx.Order("problem_basic.id DESC")
}
