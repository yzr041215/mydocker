package deamon

import (
	"engine/internal/config"
	"net"
	"os"
	"testing"
)

func Test_Start(t *testing.T) {

	path := config.Conf.EnvConf.ImagesDataDir + "/docker.sock"
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
	l, err := net.Listen("unix", path)
	if err != nil {
		t.Error(err)
		return
	}
	err = NewRouter(l).Run()
	if err != nil {
		t.Error(err)
		return
	}

}
