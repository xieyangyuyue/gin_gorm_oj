package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gin_gorm_oj/define"
)

// MyTime 是自定义的时间类型，用于处理时间的序列化和存储
type MyTime time.Time

// MarshalJSON 实现了 json.Marshaler 接口，用于将 MyTime 类型转换为 JSON 字符串
// 使用 define.DateLayout 定义的日期格式进行格式化
func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(define.DateLayout))
	return []byte(formatted), nil
}

// Value 实现了 driver.Valuer 接口，用于将 MyTime 类型转换为数据库可存储的值
// 使用 define.DateLayout 定义的日期格式进行格式化
func (t MyTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	return tTime.Format(define.DateLayout), nil
}
