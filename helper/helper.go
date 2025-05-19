package helper

import (
	"gin_gorm_oj/utils"
	"io/ioutil"
	"os"
	"time"

	"gin_gorm_oj/define"
)

// CodeSave
// 保存代码
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + utils.GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}

// CheckGoCodeValid
// 检查golang代码的合法性
func CheckGoCodeValid(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	code := string(b)
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							return false, nil
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func ToTime(num int64) time.Time {
	return time.Unix(num, 0)
}
