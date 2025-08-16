package registry

import (
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/layerdb"
	"engine/internal/pkg/repository"
	"fmt"
)

func Pull(image string) error {
	fmt.Println("-----------------pull image:", image)
	m, err := NewClient().GetManifest("library", image, "latest")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("pull manifest failed! %v", err)
	}
	// for _, manifest := range m.Manifests {
	// 	fmt.Println(manifest.Digest)
	// }
	fmt.Println("-----------------pull manifest---------------")
	fmt.Println(m.Manifests[0].Digest)
	m2, _ := NewClient().GetMinuteManifest("library", image, m.Manifests[0].Digest)
	fmt.Println("-----------------pull layers-----------------")
	ConfigDigest, err := NewClient().GetImageConfig("library", image, m2.Config.Digest)
	if err != nil {
		return fmt.Errorf("pull image config failed! %v", err)
	}
	ImageConfig, err := imagedb.GetConfig(ConfigDigest)
	if err != nil {
		return fmt.Errorf("get image config failed! %v", err)
	}
	var OnlyDiffId string
	if len(ImageConfig.RootFS.DiffIDs) == 1 {
		OnlyDiffId = ImageConfig.RootFS.DiffIDs[0]
	}

	linkpaths := make([]string, 0)
	diffDbs := make([]layerdb.DiffDb, 0)
	for _, layer := range m2.Layers {

		if d, err := NewClient().GetBlob("library", image, layer.Digest, OnlyDiffId); err != nil {
			fmt.Println("pulling layer: ", layer.Digest, " failed! ", err)
		} else {
			diffDbs = append(diffDbs, *d)
			err := d.Fun(linkpaths)
			if err != nil {
				fmt.Println(err)
			}
			linkpaths = append(linkpaths, "l/"+d.LinkId)

			fmt.Println("pulling layer: ", layer.Digest, " size: ", layer.Size, "successfully!")
		}
	}

	l := len(diffDbs)
	for i := 0; i < l; i++ {
		if i != 0 {
			diffDbs[i].SetParent(diffDbs[i-1].ChainId)
		} else {
			diffDbs[i].SetParent("")
		}
		err := diffDbs[i].Save()
		if err != nil {
			return err
		}
	}
	return repository.SaveImage(image, m2.Config.Digest, m2.Config.Digest)

}
