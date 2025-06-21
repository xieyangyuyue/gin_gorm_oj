package utils

import (
	"gin_gorm_oj/define" // 导入自定义的 define 包，其中可能包含 ValidGolangPackageMap 等定义
	"os"                 // 导入 os 包，用于操作系统相关功能，如文件和目录操作
	"time"               // 导入 time 包，用于时间相关的操作
)

// CodeSave 函数用于将用户提交的代码保存到文件系统
// code: 待保存的字节切片形式的代码内容
// 返回值: 保存成功后的文件路径和可能遇到的错误
func CodeSave(code []byte) (string, error) {
	// 构造存储代码的目录名称，使用 GetUUID() 生成唯一标识符
	dirName := "code/" + GetUUID()
	// 构造代码文件的完整路径，文件名为 main.go
	path := dirName + "/main.go"

	// 创建代码目录，权限设置为 0777（所有者、组、其他用户都可读、写、执行）
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		// 如果目录创建失败，返回空路径和错误信息
		return "", err
	}

	// 创建代码文件
	f, err := os.Create(path)
	if err != nil {
		// 如果文件创建失败，返回空路径和错误信息
		return "", err
	}

	// 将代码内容写入文件
	_, err = f.Write(code)
	if err != nil {
		return "", err
	}
	// 使用 defer 确保文件在函数返回前关闭，释放资源
	defer f.Close()

	// 返回保存成功的文件路径和 nil 错误
	return path, nil
}

// CheckGoCodeValid 函数用于检查 Golang 代码的合法性，特别是 import 语句
// 该函数旨在限制用户可以导入的包，防止恶意操作或不必要的包引入
// path: 待检查的 Golang 代码文件路径
// 返回值: 如果代码合法返回 true，否则返回 false；以及可能遇到的错误
func CheckGoCodeValid(path string) (bool, error) {
	// 读取代码文件内容
	b, err := os.ReadFile(path)
	if err != nil {
		// 如果文件读取失败，返回 false 和错误信息
		return false, err
	}

	// 将字节切片转换为字符串，以便进行字符串操作
	code := string(b)

	// 遍历代码字符串，查找 "import" 关键字
	// len(code)-6 是为了确保 "import" 后有足够的字符进行判断，避免越界
	for i := 0; i < len(code)-6; i++ {
		// 如果找到 "import" 关键字
		if code[i:i+6] == "import" {
			var flag byte // 用于存储 "import" 后的第一个非空格字符

			// 查找 "import" 后的第一个非空格字符
			for i = i + 7; i < len(code); i++ { // 从 "import" 后的第七个字符开始
				if code[i] == ' ' { // 跳过空格
					continue
				}
				flag = code[i] // 记录第一个非空格字符
				break          // 找到后退出循环
			}

			// 根据 flag 判断 import 语句的形式：是单行导入还是多行导入
			if flag == '(' { // 多行导入，例如 import ( "fmt" "os" )
				for i = i + 1; i < len(code); i++ { // 从 '(' 后的字符开始遍历
					if code[i] == ')' { // 如果遇到 ')'，说明多行导入结束
						break
					}
					if code[i] == '"' { // 如果遇到 '"'，表示开始一个包路径
						t := "" // 用于存储解析到的包路径

						// 解析包路径，直到遇到下一个 '"'
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' { // 如果遇到 '"'，说明包路径结束
								break
							}
							t += string(code[i]) // 逐个字符添加到包路径字符串
						}

						// 检查解析到的包路径是否在合法的 Golang 包映射中
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							// 如果不在，则认为代码不合法，返回 false
							return false, nil
						}
					}
				}
			} else if flag == '"' { // 单行导入，例如 import "fmt"
				t := "" // 用于存储解析到的包路径

				// 解析包路径，直到遇到下一个 '"'
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' { // 如果遇到 '"'，说明包路径结束
						break
					}
					t += string(code[i]) // 逐个字符添加到包路径字符串
				}

				// 检查解析到的包路径是否在合法的 Golang 包映射中
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					// 如果不在，则认为代码不合法，返回 false
					return false, nil
				}
			}
		}
	}
	// 如果所有 import 语句都合法，则返回 true
	return true, nil
}

// ToTime 函数将一个 Unix 时间戳（秒）转换为 time.Time 类型
// num: Unix 时间戳（int64 类型，通常表示从1970年1月1日00:00:00 UTC开始的秒数）
// 返回值: 对应的 time.Time 对象
func ToTime(num int64) time.Time {
	// time.Unix(sec, nsec) 将 Unix 时间戳转换为 time.Time
	// sec 是秒数，nsec 是纳秒数
	return time.Unix(num, 0)
}
