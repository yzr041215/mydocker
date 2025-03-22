package runc

import (
	config2 "engine/internal/config"
	"engine/internal/pkg/distribution"
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/layerdb"
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
	temp, err := layerdb.GetDiffDb(cacheid)
	OwnLink := filepath.Join(config2.Conf.EnvConf.ImagesDataDir, "overlay2", "l", temp.LinkId)

	//diffpath := filepath.Join(cachepath, "diff")
	lowerpath := filepath.Join(cachepath, "lower")
	workDir := filepath.Join(cachepath, "woke")
	lowers, err := util.ReadLowers(lowerpath)
	if err != nil {
		fmt.Println(err)
	}
	lowerOpt := ""
	for _, lower := range lowers {
		lowerOpt += lower + ":"
	}
	lowerOpt += OwnLink
	upperDir := filepath.Join(cachepath, "upper")
	mergedDir := filepath.Join(cachepath, "merged")
	util.CreateFile2(upperDir)
	util.CreateFile2(mergedDir)
	fmt.Println("lowerOpt:", lowerOpt)
	fmt.Println("upperDir:", upperDir)
	fmt.Println("mergedDir:", mergedDir)
	fmt.Println("workDir:", workDir)
	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerOpt, upperDir, workDir)
	if err := syscall.Mount("overlay", mergedDir, "overlay", 0, options); err != nil {
		fmt.Printf("挂载 OverlayFS 时出错: %v\n", err)
		return
	}

	fmt.Println("OverlayFS 挂载成功!")
}
