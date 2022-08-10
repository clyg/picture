package main

import (
	"errors"
	"fmt"
	"os"

	"strings"
)

var result = make([]string, 0)

func main() {
	dir := os.Args[1]
	r, err := searchPath(dir, IsDir, IsFull, "cos星")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	for _, v := range r {
		fmt.Println(v)
	}
}

// 函数中所需要使用的flag
const (
	IsDir     = iota // 搜索的对象是文件夹
	IsFile           // 搜索的对象是文件
	IsKeyWord        // 要进行关键字匹配
	IsFull           // 要进行全字匹配
)

// 搜索目录中的文件或者文件夹的函数
// dir为要搜索的路径
// flag为搜索文件还是文件夹
// keyflag为选择全字匹配还是部分匹配
// key为匹配的字符串
func searchPath(dir string, flag int, keyflag int, key string) ([]string, error) {
	cleanResult()
	if checkIsDir(dir) { // 如果输入的确实是一个路径并且存在
		var err error
		switch flag {
		case IsDir:
			err = scanAllDir(dir)
		case IsFile:
			err = scanAllFile(dir)
		default:
			fmt.Println("未知的文件flag")
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("输入不是一个路径")
	}
	if len(result) == 0 {
		// 如果什么都没有搜索到, 即文件夹为空
		return []string{}, nil
	}
	if (keyflag != IsKeyWord && keyflag != IsFull) || key == "" {
		// 如果给定的keyflag不正确又或者key为空则认为不进行匹配, 返回所有搜索到的文件或文件夹
		return result, nil
	} else {
		switch keyflag {
		case IsKeyWord:
			result = searchKey(key)
		case IsFull:
			result = searchFull(key)
		default:
			fmt.Println("未知的KeyFlag")
		}
	}
	return result, nil
}

// 从搜索到的字符串中将
func searchKey(keyword string) []string {
	keyDir := make([]string, 0)
	for _, v := range result {
		tmpdir := v[strings.LastIndex(v, "\\")+1:]
		if strings.Contains(tmpdir, keyword) {
			keyDir = append(keyDir, v)
		}
	}
	return keyDir
}

func searchFull(fullword string) []string {
	fullDir := make([]string, 0)
	for _, v := range result {
		tmpdir := v[strings.LastIndex(v, "\\")+1:]
		if tmpdir == fullword {
			fullDir = append(fullDir, v)
		}
	}
	return fullDir
}

func scanAllDir(dir string) error {
	if dir[len(dir)-1] == '\\' {
		dir = dir[:len(dir)-1]
	}
	result = append(result, dir)
	infos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range infos {
		if v.IsDir() {
			err := scanAllDir(dir + "\\" + v.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func scanAllFile(dir string) error {
	infos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range infos {
		if v.IsDir() {
			scanAllFile(dir + "\\" + v.Name())
		} else {
			result = append(result, dir+"\\"+v.Name())
		}
	}
	return nil
}

func cleanResult() {
	result = make([]string, 0, 0)
}

func checkIsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
