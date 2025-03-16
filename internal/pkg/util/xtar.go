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
)

var (
	i int = 0
)

// Xtar 解压TAR包
func Xtar(reader io.Reader, outpath string) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gzipReader)
	//计算sha256
	outpath = outpath + "/example" + fmt.Sprintf("%d", i)
	i++
	os.MkdirAll(outpath, 0755)
	h := NewMidSha256(tr)

	//defer func(oldpath, newpath string) {
	//	err := os.Rename(oldpath, newpath)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}(outpath, filepath.Dir(outpath)+"\\"+h.Sum())

	// 遍历TAR档案中的所有文件
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // 没有更多的文件了
		}
		if err != nil {
			return err
		}
		//fmt.Printf("Contents of %s:\n", hdr.Name)
		// 跳过目录
		if hdr.Typeflag == tar.TypeDir {
			continue
		}
		file_path := outpath + "/" + hdr.Name
		fmt.Println("Extracting", file_path)
		filedir := filepath.Dir(file_path)
		err = os.MkdirAll(filedir, 0755)
		if err != nil {
			return err
		}
		err = func() error {
			f, err := os.Create(file_path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, h); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	//fmt.Println(filepath.Dir(outpath) + "\\" + h.Sum())
	//time.Sleep(10 * time.Millisecond)
	func(oldpath, newpath string) {
		err := os.Rename(oldpath, newpath)
		if err != nil {
			fmt.Println(err)
		}
	}(outpath, filepath.Dir(outpath)+"\\"+h.Sum())
	return nil
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
