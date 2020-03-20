package schedule

import (
	"errors"
	"github.com/devfeel/rockman/state"
	"sort"
)

const (
	Balance_LowerLoad = iota //balance by lower load
	Balance_JobCount
	Balance_CpuRate
	Balance_MemoryRate
)

type (
	Scheduler struct {
	}
)

var ErrorNotSupportBalanceType = errors.New("not support balance type")
var ErrorNotExistsEnoughResource = errors.New("not exists enough resource")

// Schedule
func (s *Scheduler) Schedule(balanceType int, resources map[string]*state.ResourceInfo) ([]*state.ResourceInfo, error) {
	if resources == nil || len(resources) <= 0 {
		return nil, ErrorNotExistsEnoughResource
	}
	if balanceType == Balance_LowerLoad {
		rs := s.getSortLoadResources(resources)
		return rs, nil
	}

	if balanceType == Balance_CpuRate {
		rs := s.getSortCpuResources(resources)
		return rs, nil
	}

	if balanceType == Balance_MemoryRate {
		rs := s.getSortMemoryResources(resources)
		return rs, nil
	}

	if balanceType == Balance_JobCount {
		rs := s.getSortJobResources(resources)
		return rs, nil
	}

	return nil, ErrorNotSupportBalanceType
}

func (s *Scheduler) getSortLoadResources(resources map[string]*state.ResourceInfo) state.LoadResources {
	var loadResources state.LoadResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortCpuResources(resources map[string]*state.ResourceInfo) state.CpuResources {
	var loadResources state.CpuResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortMemoryResources(resources map[string]*state.ResourceInfo) state.MemoryResources {
	var loadResources state.MemoryResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortJobResources(resources map[string]*state.ResourceInfo) state.JobResources {
	var loadResources state.JobResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}
