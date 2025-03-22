package runc

import (
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/repository"
	"fmt"
	"testing"
)

func TestRunA(t *testing.T) {
	ImageSha256, err := repository.GetImage("nginx", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ImageSha256)
	config, err := imagedb.GetConfig(ImageSha256)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(config)

}
