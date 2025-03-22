package distribution

import (
	"engine/internal/config"
	"fmt"
	"os"
	"path/filepath"
)

func IsExistDigest(diffID string) bool {
	fileinfo, err := os.Stat(filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "diffid-by-digest", diffID))
	//如果文件存在，则直接返回
	if err == nil && fileinfo.Size() > 0 {
		fmt.Printf("Image %s already exists\n", diffID)
		return true
	}
	return false
}

// 建立digest -> diffID的映射
func SaveDiffID(digest string, diffID string) error {
	//保存diffID
	//创建diffID文件
	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "diffid-by-digest")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(path, digest))
	if err != nil {
		return err
	}
	defer f.Close()
	//写入diffID
	_, err = f.WriteString(diffID)
	if err != nil {
		return err
	}
	return nil
}
func GetDiffID(imageID string) (string, error) {
	//获取diffID
	f, err := os.Open(filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "diffid-by-digest", imageID))
	if err != nil {
		return "", err
	}
	defer f.Close()
	var diffID string
	_, err = fmt.Fscanf(f, "%s", &diffID)
	if err != nil {
		return "", err
	}
	return diffID, nil
}
