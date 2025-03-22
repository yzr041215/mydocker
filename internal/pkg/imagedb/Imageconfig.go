package imagedb

import (
	"encoding/json"
	"engine/internal/config"
	"engine/internal/pkg/util"
	"io"
	"os"
)

type DockerImageConfig struct {
	Architecture string         `json:"architecture"`
	Config       ImageConfig    `json:"config"`
	Created      string         `json:"created"`
	History      []HistoryEntry `json:"history"`
	OS           string         `json:"os"`
	RootFS       RootFS         `json:"rootfs"`
}

type ImageConfig struct {
	ExposedPorts map[string]struct{} `json:"ExposedPorts"`
	Env          []string            `json:"Env"`
	Entrypoint   []string            `json:"Entrypoint"`
	Cmd          []string            `json:"Cmd"`
	Labels       map[string]string   `json:"Labels"`
	StopSignal   string              `json:"StopSignal"`
}

type HistoryEntry struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	Comment    string `json:"comment,omitempty"`
	EmptyLayer *bool  `json:"empty_layer,omitempty"` // 使用指针处理可能缺失的布尔值
}

type RootFS struct {
	Type    string   `json:"type"`
	DiffIDs []string `json:"diff_ids"`
}

var configPath = config.Conf.EnvConf.ImagesDataDir + "/image/overlay2/imagedb/content"

func init() {
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
func GetConfig(digest string) (*DockerImageConfig, error) {
	//如果包含前缀sha256:，则去掉
	if len(digest) > 7 && digest[:7] == "sha256:" {
		digest = digest[7:]
	}
	//fmt.Println(configPath + "/" + digest)
	f, err := os.Open(configPath + "/" + digest)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := DockerImageConfig{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func SaveConfig(digest string, config *DockerImageConfig) error {
	//如果包含前缀sha256:，则去掉
	if len(digest) > 7 && digest[:7] == "sha256:" {
		digest = digest[7:]
	}
	f, err := os.Create(configPath + "/" + digest)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}
func SavaConfigByReader(reader io.Reader) error {

	var body []byte
	body, _ = io.ReadAll(reader)
	digest := util.Sha256(body)

	f, err := os.Create(configPath + "/" + digest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(body)
	if err != nil {
		return err
	}
	return f.Sync()
}
