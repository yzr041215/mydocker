package repository

import (
	"encoding/json"
	"engine/internal/config"
	"fmt"
	"os"
)

type Repositories struct {
	Repositories map[string]map[string]string `json:"repositories"`
}

var restorePath = config.Conf.EnvConf.ImagesDataDir + "/image/overlay2/repositories.json"

func init() {
	if _, err := os.Stat(restorePath); os.IsNotExist(err) {
		err := os.MkdirAll(config.Conf.EnvConf.ImagesDataDir+"/image/overlay2", 0755)
		if err != nil {
			fmt.Println("create image/overlay2 dir failed", err)
		}
		f, err := os.Create(restorePath)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString("{\"repositories\":{}}")
		if err != nil {
			panic(err)
		}
		f.Close()
	}

}
func GetImagesList() (nametags []string, err error) {
	f, err := os.Open(restorePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var repos Repositories
	err = json.NewDecoder(f).Decode(&repos)
	if err != nil {
		return nil, err
	}
	for name, tags := range repos.Repositories {
		for tag := range tags {
			nametags = append(nametags, name+":"+tag)
		}
	}
	return nametags, nil
}
func GetImage(name string, tag string) (digest string, err error) {
	f, err := os.Open(restorePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var repos Repositories
	err = json.NewDecoder(f).Decode(&repos)
	if err != nil {
		return "", err
	}
	if _, ok := repos.Repositories[name]; !ok {
		return "", fmt.Errorf("image %s not found", name)
	}
	if tag == "" { //返回第一个
		for t := range repos.Repositories[name] {
			return repos.Repositories[name][t], nil
		}
	}
	if _, ok := repos.Repositories[name][name+":"+tag]; !ok {
		return "", fmt.Errorf("tag %s not found for image %s", tag, name)
	}
	digest = repos.Repositories[name][name+":"+tag]
	return digest, nil
}
func SaveImage(name string, tag string, digest string) error {
	f, err := os.Open(restorePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var repos Repositories
	err = json.NewDecoder(f).Decode(&repos)
	if err != nil {
		return err
	}
	if _, ok := repos.Repositories[name]; !ok {
		repos.Repositories[name] = make(map[string]string)
	}
	repos.Repositories[name][name+":"+tag] = digest
	f, err = os.Create(restorePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(repos)
	if err != nil {
		return err
	}
	return nil
}
