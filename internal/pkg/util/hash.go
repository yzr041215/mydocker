package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func Sha256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的 SHA256 哈希值
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 计算目录的 SHA256 哈希值（diff_id）
func CalculateDiffID(dirPath string) (string, error) {
	var fileHashes []string

	// 遍历目录中的所有文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理普通文件
		if !info.IsDir() {
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}

			fileHash, err := CalculateFileHash(path)
			if err != nil {
				return err
			}

			// 将文件名和哈希值拼接
			fileHashes = append(fileHashes, fmt.Sprintf("%s%s", fileHash, relPath))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// 按文件名排序
	sort.Strings(fileHashes)

	// 计算整体 SHA256
	hash := sha256.New()
	for _, line := range fileHashes {
		hash.Write([]byte(line + "\n"))
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
