package registry

// ImageManifest 表示 OCI 镜像的 Manifest 文件
type ImageManifest struct {
	SchemaVersion int               `json:"schemaVersion"` // Manifest 版本（通常为 2）
	MediaType     string            `json:"mediaType"`     // Manifest 的媒体类型
	Config        ManifestConfig    `json:"config"`        // 镜像配置文件的描述
	Layers        []Layer           `json:"layers"`        // 镜像的层信息
	Annotations   map[string]string `json:"annotations"`   // 镜像的元数据注释
}

// ManifestConfig 表示镜像配置文件的描述
type ManifestConfig struct {
	MediaType string `json:"mediaType"` // 配置文件的媒体类型
	Digest    string `json:"digest"`    // 配置文件的 digest
	Size      int    `json:"size"`      // 配置文件的大小（字节）
}

// Layer 表示镜像的层信息
type Layer struct {
	MediaType string `json:"mediaType"` // 层的媒体类型
	Digest    string `json:"digest"`    // 层的 digest
	Size      int    `json:"size"`      // 层的大小（字节）
}
