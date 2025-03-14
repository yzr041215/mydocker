package test

import (
	"engine/internal/config"
	"engine/internal/pkg/registry"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(config.Conf)
	fmt.Println(registry.NewClient().GetToken("repository", "library", "redis", "pull"))
}

func TestPull(t *testing.T) {
	m, _ := registry.NewClient().GetManifest("library", "redis", "latest")
	fmt.Println("----------------------------------------")
	fmt.Println(m.Manifests[0].Digest)
	m2, _ := registry.NewClient().GetMinuteManifest("library", "redis", m.Manifests[0].Digest)
	fmt.Println("-")
	err := registry.NewClient().GetImageConfig("library", "redis", m2.Config.Digest)
	if err != nil {
		return
	}
	for _, layer := range m2.Layers {
		if err := registry.NewClient().GetBlob("library", "redis", layer.Digest); err != nil {
			fmt.Println(err)
		}
	}
}
