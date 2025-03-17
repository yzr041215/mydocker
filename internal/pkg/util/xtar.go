package util

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
)

var (
	i int = 0
)

// Xtar 解压TAR包
// func Xtar(reader io.Reader, outpath string) error {
//
//		gzipReader, err := gzip.NewReader(reader)
//		if err != nil {
//			return err
//		}
//		tr := tar.NewReader(gzipReader)
//		//计算sha256
//		outpath = outpath + "/example" + fmt.Sprintf("%d", i)
//		i++
//		os.MkdirAll(outpath, 0755)
//		h := NewMidSha256(tr)
//
//		//defer func(oldpath, newpath string) {
//		//	err := os.Rename(oldpath, newpath)
//		//	if err != nil {
//		//		fmt.Println(err)
//		//	}
//		//}(outpath, filepath.Dir(outpath)+"\\"+h.Sum())
//
//		// 遍历TAR档案中的所有文件
//		for {
//			hdr, err := tr.Next()
//			if err == io.EOF {
//				break // 没有更多的文件了
//			}
//			if err != nil {
//				return err
//			}
//			//fmt.Printf("Contents of %s:\n", hdr.Name)
//			// 跳过目录
//			if hdr.Typeflag == tar.TypeDir {
//				continue
//			}
//			file_path := outpath + "/" + hdr.Name
//			fmt.Println("Extracting", file_path)
//			filedir := filepath.Dir(file_path)
//			err = os.MkdirAll(filedir, 0755)
//			if err != nil {
//				return err
//			}
//			err = func() error {
//				f, err := os.Create(file_path)
//				if err != nil {
//					return err
//				}
//				defer f.Close()
//				if _, err := io.Copy(f, h); err != nil {
//					return err
//				}
//				return nil
//			}()
//			if err != nil {
//				return err
//			}
//		}
//		//fmt.Println(filepath.Dir(outpath) + "\\" + h.Sum())
//		//time.Sleep(10 * time.Millisecond)
//		func(oldpath, newpath string) {
//			err := os.Rename(oldpath, newpath)
//			if err != nil {
//				fmt.Println(err)
//			}
//		}(outpath, filepath.Dir(outpath)+"\\"+h.Sum())
//		return nil
//	}
//
// Xtar 解压TAR包并计算哈希值
func Xtar(reader io.Reader, outpath string) (string, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	tr := tar.NewReader(gzipReader)

	// 创建解压目录
	outpath = filepath.Join(outpath, "example")
	if err := os.MkdirAll(outpath, 0755); err != nil {
		return "", err
	}

	// 用于存储文件路径和哈希值
	fileHashes := make(map[string]string)
	hasher := sha256.New()

	// 遍历TAR档案中的所有文件
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // 没有更多的文件了
		}
		if err != nil {
			return "", err
		}

		// 处理文件路径
		filePath := filepath.Join(outpath, hdr.Name)

		// 处理不同类型的文件
		switch hdr.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg, tar.TypeRegA:
			// 创建文件
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return "", err
			}

			file, err := os.Create(filePath)
			if err != nil {
				return "", err
			}

			// 计算文件内容的哈希值
			fileHasher := sha256.New()
			if _, err := io.Copy(io.MultiWriter(file, fileHasher), tr); err != nil {
				file.Close()
				return "", err
			}
			file.Close()

			// 存储文件哈希值
			fileHashes[hdr.Name] = fmt.Sprintf("%x", fileHasher.Sum(nil))
		case tar.TypeSymlink:
			// 创建符号链接
			if err := os.Symlink(hdr.Linkname, filePath); err != nil {
				return "", err
			}
		default:
			// 忽略其他类型的文件（如设备文件）
			continue
		}
	}

	// 对文件路径进行排序
	var files []string
	for file := range fileHashes {
		files = append(files, file)
	}
	sort.Strings(files)

	// 计算最终哈希值
	hasher.Reset()
	for _, file := range files {
		hasher.Write([]byte(fileHashes[file]))
	}
	finalHash := fmt.Sprintf("%x", hasher.Sum(nil))

	// 重命名目录
	newPath := filepath.Join(filepath.Dir(outpath), finalHash)
	if err := os.Rename(outpath, newPath); err != nil {
		return "", err
	}

	return finalHash, nil
}

type MidSha256Reader struct {
	io.Reader
	H hash.Hash
}

func NewMidSha256(r io.Reader) *MidSha256Reader {
	return &MidSha256Reader{
		Reader: r,
		H:      sha256.New(),
	}
}

func (m *MidSha256Reader) Read(p []byte) (n int, err error) {
	n, err = m.Reader.Read(p)
	m.H.Write(p[:n])
	return
}

func (m *MidSha256Reader) Sum() string {
	return fmt.Sprintf("%x", m.H.Sum(nil))
}
