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
	"runtime"
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
// XGzip 解压GZIP包
func XGzip(reader io.Reader, OutFile string) (string, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()
	//计算sha256
	h := sha256.New()
	if _, err := io.Copy(h, gzipReader); err != nil {
		return "", err
	}
	file, err := os.Create(OutFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, gzipReader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Xtar 解压TAR包并计算哈希值
func Xtar(reader io.Reader, outpath string) (sha256 string, size int64, err error) {

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", -1, err
	}
	defer gzipReader.Close()
	midReader := NewMidSha256(gzipReader)
	tr := tar.NewReader(midReader)

	// 遍历TAR档案中的所有文件
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // 没有更多的文件了
		}
		if err != nil {
			return "", -1, err
		}

		// 处理文件路径
		filePath := filepath.Join(outpath, hdr.Name)

		// 处理不同类型的文件
		switch hdr.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return "", -1, err
			}
		case tar.TypeReg, tar.TypeRegA:
			// 创建文件
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return "", -1, err
			}

			file, err := os.Create(filePath)
			if err != nil {
				return "", -1, err
			}

			err = os.Chmod(filePath, hdr.FileInfo().Mode())
			if err != nil {
				file.Close()
				return "", -1, err
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return "", -1, err
			}
			file.Close()

		case tar.TypeSymlink:
			// 创建符号链接
			if runtime.GOOS == "windows" {
				// Windows不支持符号链接，忽略
				continue
			}
			if err := os.Symlink(hdr.Linkname, filePath); err != nil {
				return "", -1, err
			}
		default:
			// 忽略其他类型的文件（如设备文件）
			continue
		}
	}

	finalHash := midReader.Sum()
	return finalHash, midReader.Size(), nil
}

type MidSha256Reader struct {
	io.Reader
	H    hash.Hash
	size int64
}

func NewMidSha256(r io.Reader) *MidSha256Reader {
	return &MidSha256Reader{
		Reader: r,
		H:      sha256.New(),
		size:   0,
	}
}

func (m *MidSha256Reader) Read(p []byte) (n int, err error) {
	n, err = m.Reader.Read(p)
	m.H.Write(p[:n])
	m.size += int64(n)
	return
}

func (m *MidSha256Reader) Sum() string {
	return fmt.Sprintf("%x", m.H.Sum(nil))
}
func (m *MidSha256Reader) Size() int64 {
	return m.size
}

func ReadFull(f io.Reader) string {
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return ""
	}
	return string(buf[:n])
}
