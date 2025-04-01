package runc

import (
	config2 "engine/internal/config"
	"engine/internal/pkg/distribution"
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/repository"
	"engine/internal/pkg/util"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
)

func Mount(image string) (rootfs string, err error) {
	ImageSha256, err := repository.GetImage(image, "")
	if err != nil {
		return "", err
	}
	fmt.Println("ImageSha256:", ImageSha256)
	config, err := imagedb.GetConfig(ImageSha256)
	if err != nil {
		return "", err
	}
	//最后一个DiffID即为镜像的ID
	DiffID := config.RootFS.DiffIDs[len(config.RootFS.DiffIDs)-1]
	fmt.Println("highest DiffID:", DiffID)
	cacheid, err := distribution.GetCacheIDByDiffID(DiffID)
	if err != nil {
		return "", err
	}
	fmt.Println("cacheid:", cacheid)
	cachepath := filepath.Join(config2.Conf.EnvConf.ImagesDataDir, "overlay2", cacheid)

	OwnLink := filepath.Join(cachepath, "diff")

	//diffpath := filepath.Join(cachepath, "diff")
	lowerpath := filepath.Join(cachepath, "lower")
	workDir := filepath.Join(cachepath, "work")
	lowers, err := util.ReadLowers(lowerpath)
	if err != nil {
		fmt.Println(err)
	}
	lowerOpt := ""
	for _, lower := range lowers {
		if strings.HasSuffix(lower, "diff") {
			lowerOpt += lower + ":"
		}
	}
	lowerOpt += OwnLink
	upperDir := filepath.Join(cachepath, "upper")
	mergedDir := filepath.Join(cachepath, "merged")

	util.CreateDir(upperDir)
	util.CreateDir(mergedDir)
	util.CreateDir(workDir)
	fmt.Println("lowerOpt:", lowerOpt)
	fmt.Println("upperDir:", upperDir)
	fmt.Println("mergedDir:", mergedDir)
	fmt.Println("workDir:", workDir)

	if ok, err := util.IsMountPoint(mergedDir); ok || err != nil {
		fmt.Println("------------------have mount in mergedDir---------------")
		fmt.Println("mergedDir:", mergedDir)
		return mergedDir, err
	}
	//mount proc
	if err := util.CreateDir(filepath.Join(mergedDir, "proc")); err != nil {
		return "", fmt.Errorf("mkdir proc failed: %v", err)
	}
	if err := syscall.Mount("proc", filepath.Join(mergedDir, "proc"), "proc", 0, ""); err != nil {
		return "", fmt.Errorf("mount proc failed: %v", err)
	}
	fmt.Println("mount proc success!")

	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerOpt, upperDir, workDir)
	if err := syscall.Mount("overlay", mergedDir, "overlay", 0, options); err != nil {
		return "", fmt.Errorf("mount overlayfs failed: %v", err)
	}

	fmt.Println("OverlayFS 挂载成功!")
	return mergedDir, nil
}
