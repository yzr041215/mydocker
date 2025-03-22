package registry

import (
	"encoding/json"
	"engine/internal/config"
	"engine/internal/pkg/distribution"
	"engine/internal/pkg/imagedb"
	"engine/internal/pkg/layerdb"
	"engine/internal/pkg/util"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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

	if _, err = file.Write(body); err != nil {
		return nil, err
	}

	var manifests ImageManifest
	err = json.Unmarshal(body, &manifests)
	return &manifests, err
}
func (c *Client) GetBlob(library, image, DigestOrTag string) (d *layerdb.DiffDb, err error) {
	DigestOrTag = strings.TrimPrefix(DigestOrTag, "sha256:")
	donloadpathfile := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "overlay2")
	if library == "" {
		library = "library"
	}
	// 对应digest的diff文件是否存在 下载过就不用下载了
	if distribution.IsExistDigest(DigestOrTag) {
		//如果是第一层，直接复用 ： （digest==changeId）的时候是第一层
		if d2, err := layerdb.GetDiffDb(DigestOrTag); err == nil {
			return d2, nil
		}
		//cache_id, err := distribution.GetCacheIDByDigest(DigestOrTag)
		//if err != nil {
		//	return nil, err
		//}
		//复用layer逻辑 把新建的diff文件link到../l/cache_id/diff
		//TODO
		return nil, errors.New("exist digest " + DigestOrTag)
	}
	token, err := c.GetToken("repository", library, image, "pull")
	if err != nil {
		return nil, err
	}
	url := config.Conf.RegistryMirror + "/v2/" + library + "/" + image + "/blobs/sha256:" + DigestOrTag
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.docker.container.image.v1+json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	cache_id := util.GenerateUUID()

	diff_id, size, err := util.Xtar(res.Body, filepath.Join(donloadpathfile, cache_id, "diff"))
	if err != nil {
		return nil, err
	}

	//创建./link文件  等会链接到../l/cache_id/diff
	linkId := util.GenerateIinkUUID()
	linkpath := filepath.Join(donloadpathfile, "l")
	if err := os.MkdirAll(linkpath, 0755); err != nil {
		return nil, err
	}
	linkfile := filepath.Join(linkpath, linkId)
	//是windwos 系统的话，不创建符号链接，直接返回
	if runtime.GOOS != "windows" {
		if err := os.Symlink(filepath.Join("..", cache_id, "diff"), linkfile); err != nil {
			return nil, err
		}
	}

	//fmt.Println("diff_id: ", diff_id)
	d = layerdb.NewDiffDb(cache_id, diff_id, size, linkId)
	if err := distribution.SaveCacheID(DigestOrTag, cache_id); err != nil {
		return d, err
	}

	return d, distribution.SaveDiffID(DigestOrTag, diff_id)
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
	return imagedb.SavaConfigByReader(res.Body)

}
