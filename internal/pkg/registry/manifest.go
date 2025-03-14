package registry

// Manifest 结构体表示 Docker 镜像的 Manifest 信息
type Manifest struct {
	Annotations map[string]string `json:"annotations"` // 镜像的元数据，包含构建信息、基础镜像、架构等
	Digest      string            `json:"digest"`      // 镜像 Manifest 的唯一标识符（SHA256 哈希值）
	MediaType   string            `json:"mediaType"`   // Manifest 的媒体类型（如 application/vnd.oci.image.manifest.v1+json）
	Platform    Platform          `json:"platform"`    // 描述镜像的目标平台（架构、操作系统和变体）
	Size        int               `json:"size"`        // Manifest 文件的大小（字节）
}

// Platform 结构体描述镜像的目标平台
type Platform struct {
	Architecture string `json:"architecture"`      // 目标架构（如 amd64、arm64 等）
	OS           string `json:"os"`                // 目标操作系统（如 linux）
	Variant      string `json:"variant,omitempty"` // 架构的变体（如 v8），omitempty 表示如果为空则不包含该字段
}
type Manifests struct {
	Manifests []Manifest `json:"manifests"`
}

// GetLocalPlatform 获取本机的平台信息
func GetLocalPlatform() Platform {
	// return Platform{
	// 	Architecture: runtime.GOARCH, // 获取本机架构
	// 	OS:           runtime.GOOS,   // 获取本机操作系统
	// 	Variant:      "",             // 变体信息需要额外处理
	// }
	return Platform{
		Architecture: "amd64",
		OS:           "linux",
		Variant:      "",
	}
}
func FliterManifests(manifests *Manifests) *Manifest {
	platform := GetLocalPlatform()
	for _, manifest := range manifests.Manifests {
		if manifest.Platform.Architecture == platform.Architecture && manifest.Platform.OS == platform.OS && manifest.Platform.Variant == platform.Variant {
			return &manifest
		}
	}
	return nil
}
