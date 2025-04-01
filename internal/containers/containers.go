package containers

import (
	"encoding/json"
	"engine/internal/config"
	"os"
	"path/filepath"
)

type ContainerInfo struct {
	ContainerID string `json:"ContainerID"`
	Image       string `json:"Image"`
	Pid         int    `json:"Pid"`
	Command     string `json:"Command"`
	Created     string `json:"Created"`
	Status      string `json:"Status"`
}

var (
	path = config.Conf.EnvConf.ImagesDataDir + "/containers"
)

func init() {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}
func GetContainerInfo(containerID string) (*ContainerInfo, error) {
	path1 := filepath.Join(path, containerID, "config.v2.json")
	if f, err := os.Open(path1); err != nil {
		return nil, err
	} else {
		defer f.Close()
		var info ContainerInfo
		if err := json.NewDecoder(f).Decode(&info); err != nil {
			return nil, err
		}
		return &info, nil
	}

}
func SaveContainerInfo(info *ContainerInfo) error {
	err := os.MkdirAll(path+"/"+info.ContainerID, os.ModePerm)
	if err != nil {
		return err
	}
	path1 := filepath.Join(path, info.ContainerID, "config.v2.json")
	logepath := filepath.Join(path, info.ContainerID, "log.log")
	if f, err := os.Create(path1); err != nil {
		return err
	} else {
		defer f.Close()
		if err := json.NewEncoder(f).Encode(info); err != nil {
			return err
		}
	}
	if f, err := os.Create(logepath); err != nil {
		return err
	} else {
		defer f.Close()
		//TODO: write log to file
	}
	return nil
}
