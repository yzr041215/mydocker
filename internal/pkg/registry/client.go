package registry

import (
	"encoding/json"
	"engine/internal/config"
	"engine/internal/pkg/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient, //设置默认的http客户端
	}
}

// <type>:<namespace>/<repository>:<action>
func (c *Client) GetToken(Resourcetype, namespace, repository, action string) (string, error) {
	//curl "https://docker.hlmirror.com/token?service=registry.docker.io&scope=repository:library/redis:pull"
	//如果repository中不包含namespace，则添加默认的library
	if namespace == "" {
		namespace = "library"
	}
	url := config.Conf.RegistryMirror + "/token?service=registry.docker.io&scope=" + "repository:" + namespace + "/" + repository + ":" + action
	res, err := c.client.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var token struct {
		Token string `json:"token"`
	}
	//	fmt.Println(string(body))
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
func (c *Client) GetManifest(library, image, tag string) (*Manifests, error) {
	if library == "" {
		library = "library"
	}
	if tag == "" {
		tag = "latest"
	}
	token, err := c.GetToken("repository", library, image, "pull")
	if err != nil {
		return nil, err
	}

	url := config.Conf.RegistryMirror + "/v2/" + library + "/" + image + "/manifests/" + tag

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var manifests Manifests
	//fmt.Println("tag: ", tag)
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &manifests)
	if err != nil {
		return nil, err
	}
	return &manifests, nil
}
func (c *Client) GetMinuteManifest(library, image, DigestOrTag string) (*ImageManifest, error) {
	if library == "" {
		library = "library"
	}
	if DigestOrTag == "" {
		DigestOrTag = "latest"
	}
	token, err := c.GetToken("repository", library, image, "pull")
	if err != nil {
		return nil, err
	}
	url := config.Conf.RegistryMirror + "/v2/" + library + "/" + image + "/manifests/" + DigestOrTag
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	filepath := config.Conf.EnvConf.ImagesDataDir + "/" + "manifest" + ".json"
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		return nil, err
	}
	var manifests ImageManifest
	err = json.Unmarshal(body, &manifests)
	if err != nil {
		return nil, err
	}
	return &manifests, nil
}
func (c *Client) GetBlob(library, image, DigestOrTag string) error {
	DigestOrTag = strings.TrimPrefix(DigestOrTag, "sha256:")
	donloadpathfile := config.Conf.EnvConf.ImagesDataDir + "/" + "layer"
	if library == "" {
		library = "library"
	}
	token, err := c.GetToken("repository", library, image, "pull")
	if err != nil {
		return err
	}
	url := config.Conf.RegistryMirror + "/v2/" + library + "/" + image + "/blobs/sha256:" + DigestOrTag
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.docker.container.image.v1+json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Println("Get Blob" + DigestOrTag)
	_, err = util.Xtar(res.Body, donloadpathfile)

	return err
}
func (c *Client) GetImageConfig(library, image, DigestOrTag string) error {
	url := config.Conf.RegistryMirror + "/v2/" + library + "/" + image + "/blobs/" + DigestOrTag
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	token, err := c.GetToken("repository", library, image, "pull")
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.container.image.v1+json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var body []byte
	body, _ = io.ReadAll(res.Body)
	//fmt.Println(string(body))
	filepath := config.Conf.EnvConf.ImagesDataDir + "/" + "config" + ".json"
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		return err
	}
	return nil
}
