package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func FindAllDir(dir string) []string {
	dirSlice := make([]string, 0)
	dirSlice = append(dirSlice, dir)
	dirChan := make(chan string, 10)
	tokenChan := make(chan struct{}, 32)
	var n sync.WaitGroup
	n.Add(1)
	go func() {
		walkDir(dir, &n, dirChan, tokenChan)
	}()
	go func() {
		n.Wait()
		close(dirChan)
	}()
	for dir := range dirChan {
		dirSlice = append(dirSlice, dir)
	}
	return dirSlice
}

func walkDir(dir string, n *sync.WaitGroup, dirchan chan<- string, tokenChan chan struct{}) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			tokenChan <- struct{}{}
			subdir := filepath.Join(dir, entry.Name())
			dirchan <- subdir
			go walkDir(subdir, n, dirchan, tokenChan)
			<-tokenChan
		}
	}
}

func FindAllFile(dir string) []string {
	fileSlice := make([]string, 0)
	fileChan := make(chan string, 32)
	tokenChan := make(chan struct{}, 32)
	var n sync.WaitGroup
	n.Add(1)
	go func() {
		walkFile(dir, &n, fileChan, tokenChan)
	}()
	go func() {
		n.Wait()
		close(fileChan)
	}()
	for fileName := range fileChan {
		fileSlice = append(fileSlice, fileName)
	}
	return fileSlice
}

func walkFile(dir string, n *sync.WaitGroup, fileChan chan<- string, tokenChan chan struct{}) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		subDir := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			n.Add(1)
			go walkFile(subDir, n, fileChan, tokenChan)
		} else {
			tokenChan <- struct{}{}
			fileChan <- subDir
			<-tokenChan
		}
	}
}

func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取目录下的文件信息错误", err)
		return nil
	}
	return entries
}
