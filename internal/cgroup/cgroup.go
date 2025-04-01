package cgroup

import (
	"engine/internal/cgroup/resource"
	"engine/internal/cgroup/subsystem"
	"engine/internal/pkg/util"
	"log"
)

type CgroupManager interface {
	Apply(pid int) error
	Set(res *resource.ResourceConfig) error
	Destroy() error
}

func NewCgroupManager(path string) CgroupManager {
	if util.IsCgroup2UnifiedMode() {
		log.Println("use cgroup v2")
		return NewCgroupManagerV2(path, subsystem.Subsystems)
	}
	log.Println("use cgroup v1")
	return NewCgroupManagerV1(path, subsystem.Subsystems)
}
