package tools

import (
	"os"
	"path/filepath"
	"strings"
)

func IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func DeleteFile(s string) {
	_ = os.Remove(s)
}

func CreateDirIfNotExists(filePathStr string) {
	absPath, _ := filepath.Abs(filePathStr)
	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		// 创建文件夹
		_ = os.MkdirAll(absPath, os.ModePerm)
	}
}
func UpperString(tableName string) string {
	r := camelString(tableName)
	if len(r) > 0 {
		r = strings.ToUpper(r[:1]) + r[1:]
	}
	return r
}
func LowerString(tableName string) string {
	r := camelString(tableName)
	if len(r) > 0 {
		r = strings.ToLower(r[:1]) + r[1:]
	}
	return r
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}
