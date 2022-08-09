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

const (
	IsDir = iota
	IsFile
	IsKey
	IsFull
)

func searchPath(dir string, flag int, keyflag int, key string) ([]string, error) {
	cleanResult()
	if checkIsDir(dir) {
		switch flag {
		case IsDir:
			scanAllDir(dir)
			result = keyswitch(keyflag, key)
		case IsFile:
			scanAllFile(dir)
			result = keyswitch(keyflag, key)
		default:
			fmt.Println("未知的flag")
		}
	} else {
		return nil, errors.New("输入不是一个路径")
	}

	return result, nil
}

func keyswitch(keyflag int, word string) []string {
	keyDirs := make([]string, 0)
	switch keyflag {
	case -1:
		return result
	case IsKey:
		keyDirs = searchKey(word)
	case IsFull:
		keyDirs = searchFull(word)
	default:
		fmt.Println("未知！")
	}

	return keyDirs
}

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
