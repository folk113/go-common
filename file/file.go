package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func makeDir(name string) bool {
	path := filepath.Dir(name)
	// 判断文件路径是否存在
	_,err := os.Stat(path)
	if err == nil {
		return true
	}

	err = os.MkdirAll(path, 0644)
	if err == nil {
		return true
	} else {
		fmt.Println("create file dir:", err)
		return false
	}
}

func Write(name, content string) {
	createSuccess := makeDir(name)
	if !createSuccess {
		return
	}
	file, err := os.OpenFile(name, os.O_CREATE, 0644)
	if err == nil {
		file.WriteString(content)
		fmt.Println("write file success:", name)
	} else {
		fmt.Println("write file failed:", err)
	}
}

func Read(name string) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("read file failed", name)
		return "", errors.New("read file failed")
	} else {
		return string(data), nil
	}
}