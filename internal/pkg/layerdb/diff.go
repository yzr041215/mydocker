package layerdb

import (
	"engine/internal/config"
	"engine/internal/pkg/util"
	"os"
	"path/filepath"
	"strconv"
)

type DiffDb struct {
	CacheId string `json:"cache_id"`
	DiffId  string `json:"diff_id"`
	Size    int64  `json:"size"`
	Parent  string `json:"parent"` // parent chain of this diff
	ChainId string `json:"chain_id"`
	LinkId  string `json:"-"`
	Fun     func([]string) error
}

func NewDiffDb(cacheId string, diffId string, size int64, linkid string) *DiffDb {
	return &DiffDb{
		CacheId: cacheId,
		DiffId:  diffId,
		Size:    size,
		LinkId:  linkid,
		Fun: func(lower []string) error {
			path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "overlay2", cacheId)
			os.MkdirAll(path, os.ModePerm)
			path = filepath.Join(path, "lower")
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()
			for _, l := range lower {
				if _, err := f.WriteString(l + ":"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
func (d *DiffDb) SetParent(parent string) (chainId string) {

	d.Parent = parent
	if parent == "" {
		d.ChainId = d.DiffId
		return d.DiffId
	}
	d.ChainId = util.Sha256([]byte(parent + " " + d.DiffId))
	return d.ChainId
}
func (d *DiffDb) Save() error {
	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "layerdb", "sha256")

	path = filepath.Join(path, d.ChainId)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	//7e83a2aa65f455620126a3e731fac20a1ae130eb505af09d30deabd1c4e9663d
	if f, err := os.Create(filepath.Join(path, "size")); err == nil {
		defer f.Close()
		if _, err := f.WriteString(strconv.FormatInt(d.Size, 10)); err != nil {
			return err
		}
	} else {
		return err
	}

	if f, err := os.Create(filepath.Join(path, "cache-id")); err == nil {
		defer f.Close()
		if _, err := f.WriteString(d.CacheId); err != nil {
			return err
		}
	} else {
		return err
	}

	if f, err := os.Create(filepath.Join(path, "diff")); err == nil {
		defer f.Close()
		if _, err := f.WriteString(d.DiffId); err != nil {
			return err
		}
	} else {
		return err
	}

	if f, err := os.Create(filepath.Join(path, "parent")); err == nil {
		defer f.Close()
		if _, err := f.WriteString(d.Parent); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func GetDiffDb(ChainId string) (*DiffDb, error) {
	path := filepath.Join(config.Conf.EnvConf.ImagesDataDir, "image", "overlay2", "layerdb", "sha256")

	path = filepath.Join(path, ChainId)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	diffDb := &DiffDb{}
	if f, err := os.Open(filepath.Join(path, "size")); err == nil {
		defer f.Close()
		if size, err := strconv.ParseInt(util.ReadFull(f), 10, 64); err == nil {
			diffDb.Size = size
		}
	} else {
		return nil, err
	}

	if f, err := os.Open(filepath.Join(path, "cache-id")); err == nil {
		defer f.Close()
		diffDb.CacheId = util.ReadFull(f)
	} else {
		return nil, err
	}

	if f, err := os.Open(filepath.Join(path, "diff")); err == nil {
		defer f.Close()
		diffDb.DiffId = util.ReadFull(f)
	} else {
		return nil, err
	}

	if f, err := os.Open(filepath.Join(path, "parent")); err == nil {
		defer f.Close()
		diffDb.Parent = util.ReadFull(f)
	} else {
		return nil, err
	}

	return diffDb, nil
}
