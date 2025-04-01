package subsystem

import (
	"engine/internal/cgroup/resource"
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuSubsystem struct{}

func (s *CpuSubsystem) Name() string {
	return "cpu"
}

func (s *CpuSubsystem) Set(cgroupPath string, config *resource.ResourceConfig) error {
	if config.CpuCfsQuota != 0 {
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		// 这里只是简单的计算了下，并没有处理一些特殊情况，比如负数什么的
		if err := os.WriteFile(path.Join(cgroupPath, "cpu.max"), []byte(fmt.Sprintf("%s %s", strconv.Itoa(100000/100*config.CpuCfsQuota), strconv.Itoa(100000))), 0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}
	return nil
}
