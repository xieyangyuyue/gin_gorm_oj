package utils

import uuid "github.com/satori/go.uuid" // 导入 go.uuid 库，并将其别名为 uuid，用于生成通用唯一标识符（UUID）

// GetUUID 函数用于生成一个版本4的UUID字符串。
// UUID（Universally Unique Identifier）是一个128位的数字，用于唯一标识信息。
// 版本4的UUID是基于随机数或伪随机数生成的。
func GetUUID() string { // 定义一个名为 GetUUID 的公共函数，它不接受任何参数，并返回一个字符串类型的值
	return uuid.NewV4().String() // 调用 uuid 库的 NewV4() 方法生成一个版本4的UUID对象，然后调用 String() 方法将其转换为字符串格式并返回
}
