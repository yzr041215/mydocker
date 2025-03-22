package distribution

import (
	"engine/internal/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 建立digest -> cacheID的映射
func SaveCacheID(digest string, CacheID string) error {

	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "cacheid-by-digest")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(path, digest))
	if err != nil {
		return err
	}
	defer f.Close()
	//写入CacheID
	_, err = f.WriteString(CacheID)
	if err != nil {
		return err
	}
	return nil
}

func GetCacheIDByDigest(imageID string) (string, error) {
	//去除imageID前缀
	imageID = strings.TrimPrefix(imageID, "sha256:")
	f, err := os.Open(filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "cacheid-by-digest", imageID))
	if err != nil {
		return "", err
	}
	defer f.Close()
	var CacheID string
	_, err = fmt.Fscanf(f, "%s", &CacheID)
	if err != nil {
		return "", err
	}
	return CacheID, nil
}
