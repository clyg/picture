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

func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取目录下的文件信息错误", err)
		return nil
	}
	return entries
}
