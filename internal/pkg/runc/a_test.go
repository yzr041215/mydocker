package runc

import (
	config2 "engine/internal/config"
	"engine/internal/pkg/distribution"
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/repository"
	"engine/internal/pkg/util"
	"fmt"
	"path/filepath"
	"testing"
)

func TestRunA(t *testing.T) {
	//ss
	ImageSha256, err := repository.GetImage("nginx", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("ImageSha256:", ImageSha256)
	config, err := imagedb.GetConfig(ImageSha256)
	if err != nil {
		t.Error(err)
		return
	}
	//最后一个DiffID即为镜像的ID
	DiffID := config.RootFS.DiffIDs[len(config.RootFS.DiffIDs)-1]
	fmt.Println("highest DiffID:", DiffID)
	cacheid, err := distribution.GetCacheIDByDiffID(DiffID)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("cacheid:", cacheid)
	cachepath := filepath.Join(config2.Conf.EnvConf.ImagesDataDir, "overlay2", cacheid)
	diffpath := filepath.Join(cachepath, "diff")
	lowerpath := filepath.Join(cachepath, "lower")
	wokepath := filepath.Join(cachepath, "woke")
	lowers, err := util.ReadLowers(lowerpath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("links:", lowers)
	fmt.Println("diffpath:", diffpath)
	fmt.Println("wokepath:", wokepath)
}
