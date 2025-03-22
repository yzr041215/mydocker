package distribution

import (
	"encoding/json"
	"engine/internal/config"
	"os"
	"path/filepath"
)

//[{"Digest":"sha256:d6c5e428cfd7eb02569d3dba6ac00aaa96e7e5374714920663a54645b20b7b5d","SourceRepository":"docker.io/library/redis","HMAC":""}]

type Metadata struct {
	Digest           string `json:"Digest"`
	SourceRepository string `json:"SourceRepository"`
	HMAC             string `json:"HMAC"`
}

func GetMetadataByDiffId(diffId string) *Metadata {
	FilePath := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "v2metadata-by-diffid", diffId)
	f, err := os.Open(FilePath)
	if err != nil {
		return nil
	}
	defer f.Close()
	var metadata Metadata
	if err := json.NewDecoder(f).Decode(&metadata); err != nil {
		return nil
	}
	return &metadata
}

func SaveMetadata(metadata *Metadata) error {
	FilePath := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "distribution", "v2metadata-by-diffid", metadata.Digest)
	f, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(metadata); err != nil {
		return err
	}
	return nil
}
