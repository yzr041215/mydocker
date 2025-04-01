package cgroup

import (
	"engine/internal/cgroup/resource"
	"engine/internal/cgroup/subsystem"
)

type CgroupManagerV1 struct {
}

func NewCgroupManagerV1(path string, subsystems map[string]subsystem.Subsystem) *CgroupManagerV2 {
	return &CgroupManagerV2{
		path:       path,
		subsystems: subsystems,
	}

}

func (this *CgroupManagerV1) Apply(pid int) error {
	//TODO: implement this function
	return nil
}

func (this *CgroupManagerV1) Set(res *resource.ResourceConfig) error {
	//TODO: implement this function
	return nil
}
func (this *CgroupManagerV1) Destroy() error {
	//TODO: implement this function
	return nil
}
