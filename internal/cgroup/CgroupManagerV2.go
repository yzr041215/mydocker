package cgroup

import (
	"engine/internal/cgroup/resource"
	"engine/internal/cgroup/subsystem"
	"os"
)

type CgroupManagerV2 struct {
	path       string
	subsystems map[string]subsystem.Subsystem
}

func NewCgroupManagerV2(path string, subsystems map[string]subsystem.Subsystem) *CgroupManagerV2 {
	return &CgroupManagerV2{
		path:       path,
		subsystems: subsystems,
	}

}

func (this *CgroupManagerV2) Apply(pid int) error {
	path := this.path + "/" + "mydocker"
	return os.WriteFile(path+"/tasks", []byte(string(rune(pid))), 0644)
}

func (this *CgroupManagerV2) Set(res *resource.ResourceConfig) error {
	for _, sub := range this.subsystems {
		if err := sub.Set("sys/fs/cgroup/mydocker", res); err != nil {
			return err
		}
	}
	return nil
}
func (this *CgroupManagerV2) Destroy() error {
	//TODO: implement this function
	return nil
}
