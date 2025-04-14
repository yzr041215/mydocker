package util

import (
	config2 "engine/internal/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func FormatTimeAgo(stringTime string) string {
	t, _ := strconv.ParseInt(stringTime, 10, 64)
	duration := time.Since(time.Unix(t, 0))

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%2d mins  ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%2d hours ago", int(duration.Hours()))
	default:
		return fmt.Sprintf("%2d days  ago", int(duration.Hours()/24))
	}
}
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
	//fmt.Println("temp:", temp)
	if len(temp) == 0 {
		return []string{}, nil
	}
	links = strings.Split(string(temp), ":")
	for i := 0; i < len(links); i++ {
		links[i] = config2.Conf.EnvConf.ImagesDataDir + "/overlay2/" + links[i]
	}
	fmt.Println("links:", links)
	return links, nil
}

func CreateFile1(path string) (f *os.File, err error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func CreateFile2(path string) (err error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
func CreateDir(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func IsMountPoint(path string) (bool, error) {
	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	// 获取该路径的文件信息
	var stat syscall.Stat_t
	if err := syscall.Stat(absPath, &stat); err != nil {
		return false, err
	}
	dev := stat.Dev

	// 获取父目录的文件信息
	parent := filepath.Dir(absPath)
	if err := syscall.Stat(parent, &stat); err != nil {
		return false, err
	}
	parentDev := stat.Dev

	// 如果设备号不同，则说明是挂载点
	return dev != parentDev, nil
}

const (
	unifiedMountpoint = "/sys/fs/cgroup"
)

var (
	isUnifiedOnce sync.Once
	isUnified     bool
)

func IsCgroup2UnifiedMode() bool {
	isUnifiedOnce.Do(func() {
		var st unix.Statfs_t
		err := unix.Statfs(unifiedMountpoint, &st)
		if err != nil && os.IsNotExist(err) {
			// For rootless containers, sweep it under the rug.
			isUnified = false
			return
		}
		isUnified = st.Type == unix.CGROUP2_SUPER_MAGIC
	})
	return isUnified
}
