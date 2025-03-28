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
	configpath := "/home/yzr/Documents/mydocker/config.json"
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
