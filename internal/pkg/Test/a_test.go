package test

import (
	"archive/tar"
	"compress/gzip"
	"engine/internal/config"
	"engine/internal/pkg/registry"
	"engine/internal/pkg/util"
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(config.Conf)
	fmt.Println(registry.NewClient().GetToken("repository", "library", "redis", "pull"))
}

func TestPull(t *testing.T) {
	m, err := registry.NewClient().GetManifest("library", "nginx", "latest")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("------------------------------------------")
	fmt.Println(m.Manifests[0].Digest)
	m2, _ := registry.NewClient().GetMinuteManifest("library", "nginx", m.Manifests[0].Digest)
	fmt.Println("-")
	err = registry.NewClient().GetImageConfig("library", "nginx", m2.Config.Digest)
	if err != nil {
		return
	}
	for _, layer := range m2.Layers {
		fmt.Println("GetyBlob:", layer.Digest)
		if err := registry.NewClient().GetBlob("library", "nginx", layer.Digest); err != nil {
			fmt.Println(err)
		}
	}
}

func TestHash(t *testing.T) {
	// 解压后的目录路径as
	dirPath := "D:\\VS_Code_Project\\go\\MyDocker\\data\\layer"

	f, _ := os.Open("D:\\VS_Code_Project\\go\\MyDocker\\data\\7cf63256a31a4cc44f6defe8e1af95363aee5fa75f30a248d95cae684f87c53c")
	err := util.Xtar(f, dirPath)
	if err != nil {
		fmt.Println(err)
	}
}
func TestXtar(t *testing.T) {
	infile := "D:\\VS_Code_Project\\go\\MyDocker\\data\\7cf63256a31a4cc44f6defe8e1af95363aee5fa75f30a248d95cae684f87c53c"
	outfile := "D:\\VS_Code_Project\\go\\MyDocker\\data\\AA"
	os.Mkdir(outfile, os.ModePerm)
	inFile, _ := os.Open(infile)

	defer inFile.Close()
	// 使用 gzip 解压缩
	gzipReader, err := gzip.NewReader(inFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gzipReader.Close()
	tr := tar.NewReader(gzipReader)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		fmt.Println(hdr.Name)
		os.MkdirAll(path.Join(outfile, hdr.Name), os.ModePerm)
		outFile, _ := os.Create(path.Join(outfile, hdr.Name))
		defer outFile.Close()
		io.Copy(outFile, tr)
		if err != nil {
			fmt.Println(err)
		}
	}

}
