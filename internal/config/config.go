package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	RegistryMirror string    `json:"registry-mirror"` // 镜像仓库地址列表
	EnvConf        EnvConfig `json:"env"`
}
type EnvConfig struct {
	ImagesDataDir string `json:"image_data_dir"` //  镜像数据目录
}

var (
	Conf *Config
)

func init() {
	Conf = &Config{
		RegistryMirror: "https://docker.hlmirror.com", // 镜像仓库地址列表",
		EnvConf: EnvConfig{
			ImagesDataDir: "/home/yzr/mydocker",
		},
	}

	configpath := "/home/yzr/mydocker/config.json"
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		return
	}
	Conf = &Config{}
	f, err := os.Open(configpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(Conf)
	if err != nil {
		panic(err)
	}
}
