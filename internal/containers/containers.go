package containers

import (
	"encoding/json"
	"engine/internal/config"
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

type ContainerInfo struct {
	ContainerID string   `json:"ContainerID"`
	Image       string   `json:"Image"`
	Pid         int      `json:"Pid"`
	Command     string   `json:"Command"`
	Created     string   `json:"Created"`
	Status      string   `json:"Status"`
	Volume      []string `json:"volume"`      // 容器挂载的 volume
	NetworkName string   `json:"networkName"` // 容器所在的网络
	PortMapping []string `json:"portmapping"` // 端口映射
	IP          string   `json:"ip"`
}

var (
	path = config.Conf.EnvConf.ImagesDataDir + "/containers"
)

func init() {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}
func StopContainer(containerID string) error {
	info, err := GetContainerInfo(containerID)
	if err != nil {
		return err
	}
	if err := syscall.Kill(info.Pid, syscall.SIGTERM); err != nil {
		return err
	}

	info.Pid = 0
	info.Status = "exited"
	return SaveContainerInfo(info)

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
	//logepath := filepath.Join(path, info.ContainerID, "log.log")
	if f, err := os.Create(path1); err != nil {
		return err
	} else {
		defer f.Close()
		if err := json.NewEncoder(f).Encode(info); err != nil {
			return err
		}
	}

	return nil
}
func RemoveContainerInfo(containerID string) error {

	path1 := filepath.Join(path, containerID)
	if _, err := os.Stat(path1); os.IsNotExist(err) {
		return errors.New("container not exist")
	}
	if err := os.RemoveAll(path1); err != nil {
		return err
	}
	return nil
}
func UpdateAllContainerStatus() (map[string]*ContainerInfo, error) {
	Infos := make(map[string]*ContainerInfo)
	files, err := filepath.Glob(path + "/*/config.v2.json")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		containerID := filepath.Base(filepath.Dir(file))
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		var info ContainerInfo
		if err := json.NewDecoder(f).Decode(&info); err != nil {
			return nil, err
		}
		pid := info.Pid
		if !isPidExist(pid) {
			info.Status = "exited"
		} else {
			info.Status = "running"
		}
		//refresh info
		SaveContainerInfo(&info)
		Infos[containerID] = &info
	}
	return Infos, nil
}
func isPidExist(pid int) bool {
	// 发送信号0到进程，如果进程存在且有权访问则返回nil
	err := syscall.Kill(pid, 0)
	return err == nil
}

func GetLog(containerID string) (string, error) {
	path1 := filepath.Join(path, containerID, "log.log")
	if f, err := os.Open(path1); err != nil {
		return "", err
	} else {
		defer f.Close()
		res := ""
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err != nil {
				break
			}
			res += string(buf[:n])
		}
		return res, nil
	}
}
