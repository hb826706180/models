package models

import (
	"fmt"
	"os"
	"path"
)

type File_ struct {
}

// CreateFile 创建文件
func (File_) CreateFile(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			fmt.Println("创建文件失败！", err)
			return false
		}
	}
	return true
}

// CreatePath 创建路径
func (File_) CreatePath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println("创建路径失败！", err)
			return err
		}
	}
	return nil
}

// GetPath 获取路径
func (File_) GetPath(filePath string) string {
	return path.Dir(filePath)
}

// GetName  获取文件名
func (File_) GetName(filePath string) string {
	return path.Base(filePath)
}

// 读取文件
func File读取文件全部内容(路径 string) (string, error) {
	file, err := os.ReadFile(路径)
	if err != nil {
		fmt.Println("打开文件失败>err", err)
		return "", err
	}
	return string(file), nil
}

// IfExist 判断文件或路径是否存在
func (File_) IfExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

// DelPath 删除路径和文件
func (File_) DelPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("错误 RemoveAll():", err)
		return
	}
}
