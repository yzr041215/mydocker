package util

import (
	"io"
	"os"
	"strings"
)

func ReadLowers(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	links := []string{}
	//读取所有，使用":"分割
	temp, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(temp) > 0 && temp[len(temp)-1] == ':' {
		//去掉最后一个":"
		temp = temp[:len(temp)-1]
	}
	links = strings.Split(string(temp), ":")

	return links, nil
}
