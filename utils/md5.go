package utils // 声明包名为 utils，通常用于存放工具函数，以便在其他包中引用

import (
	"crypto/md5" // 导入 Go 语言标准库的 crypto/md5 包，用于计算MD5哈希值
	"fmt"        // 导入 Go 语言标准库的 fmt 包，用于格式化输出
)

// GetMd5 函数用于计算给定字符串的MD5哈希值。
// MD5（Message-Digest Algorithm 5）是一种广泛使用的密码散列函数，可以生成一个128位（16字节）的哈希值。
// 它通常用于验证数据完整性，但已不推荐用于密码存储或数字签名等安全性要求高的场景。
//
// 参数:
//
//	s string: 待计算MD5哈希值的输入字符串。
//
// 返回值:
//
//	string: 返回输入字符串的MD5哈希值的十六进制字符串表示。
func GetMd5(s string) string {
	// md5.Sum([]byte(s)) 计算输入字符串 s 的MD5哈希值。
	// 它接收一个字节切片作为输入，并返回一个 [16]byte 类型的数组，表示128位的MD5哈希值。
	hash := md5.Sum([]byte(s))

	// fmt.Sprintf("%x", hash) 将16字节的哈希数组格式化为十六进制字符串。
	// %x 格式化动词会将字节数组中的每个字节转换为两位十六进制数。
	return fmt.Sprintf("%x", hash)
}
