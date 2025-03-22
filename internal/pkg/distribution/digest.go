package distribution

import (
	"engine/internal/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 建立diffID -> digest的映射
func SaveDigest(digest string, diffID string) error {
	//保存diffID
	//创建diffID文件
	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "digest-by-diffid")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(path, diffID))
	if err != nil {
		return err
	}
	defer f.Close()
	//写入digest
	_, err = f.WriteString(digest)
	if err != nil {
		return err
	}
	return nil
}
func GetCacheIDByDiffID(diffID string) (string, error) {
	diffID = strings.TrimPrefix(diffID, "sha256:")
	Digest, err := GetDigestByDiffID(diffID)
	if err != nil {
		return "", err
	}
	return GetCacheIDByDigest(Digest)
}

func GetDigestByDiffID(imageID string) (string, error) {
	imageID = strings.TrimPrefix(imageID, "sha256:")
	//获取digest
	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "digest-by-diffid", imageID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("diffID %s not found", imageID)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	//读取digest
	var digest string
	_, err = fmt.Fscanf(f, "%s", &digest)
	if err != nil {
		return "", err
	}
	return digest, nil
}
