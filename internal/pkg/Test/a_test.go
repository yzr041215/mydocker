package test

import (
	"engine/internal/config"
	"engine/internal/pkg/layerdb"
	"engine/internal/pkg/registry"
	"engine/internal/pkg/util"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(config.Conf)
	fmt.Println(registry.NewClient().GetToken("repository", "library", "redis", "pull"))
}
func TestPull2(t *testing.T) {
	if err := registry.Pull("nginx"); err != nil {
		fmt.Println(err)
	}
}
func TestPull(t *testing.T) {
	image := "mysql"
	m, err := registry.NewClient().GetManifest("library", image, "latest")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("--------------------------------------------")
	fmt.Println(m.Manifests[0].Digest)
	m2, _ := registry.NewClient().GetMinuteManifest("library", image, m.Manifests[0].Digest)
	fmt.Println("-")
	err = registry.NewClient().GetImageConfig("library", image, m2.Config.Digest)
	if err != nil {
		fmt.Println(err)
		return
	}
	diffDbs := make([]layerdb.DiffDb, 0)
	for _, layer := range m2.Layers {
		fmt.Println("GetyBlob: ", layer.Digest)

		if d, err := registry.NewClient().GetBlob("library", image, layer.Digest); err != nil {
			fmt.Println(err)
		} else {

			diffDbs = append(diffDbs, *d)
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
			return
		}
	}
}

func TestUUID(t *testing.T) {
	Id := util.GenerateUUID()
	fmt.Println(Id)
	Id2 := util.GenerateIinkUUID()
	fmt.Println(Id2)
}
func TestHash(t *testing.T) {
	src := "2f0051eaf9cdda36dbc23ac85e719c50a8185143bcd9bf5092f1dad6eb5d3772 2afafa4063fa83a4580946419bcb17fd7c1109691b9818c75a0893d9dbbbe2f2"
	s := util.Sha256([]byte(src))
	fmt.Println(s)
}
