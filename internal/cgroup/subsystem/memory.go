package subsystem

import (
	"engine/internal/cgroup/resource"
	"fmt"
	"os"
	"path"
)

type MemorySubsystem struct{}

func (m *MemorySubsystem) Name() string {
	return "memory"
}
func (m *MemorySubsystem) Set(cgroupPath string, config *resource.ResourceConfig) error {
	if config.MemoryLimit == "" {
		return nil
	}

	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err := os.WriteFile(path.Join(cgroupPath, "memory.max"), []byte(config.MemoryLimit), 0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	return nil
}
