package scheduler

import (
	"errors"
	"sort"
	"sync"
)

const (
	Balance_LowerLoad = iota //balance by lower load
	Balance_JobCount
	Balance_CpuRate
	Balance_MemoryRate
)

type (
	Scheduler struct {
		resources      map[string]*ResourceInfo
		resourceLocker *sync.RWMutex
	}
)

var ErrorNotSupportBalanceType = errors.New("not support balance type")
var ErrorNotExistsEnoughResource = errors.New("not exists enough resource")

func NewScheduler() *Scheduler {
	scheduler := new(Scheduler)
	scheduler.resources = make(map[string]*ResourceInfo)
	scheduler.resourceLocker = new(sync.RWMutex)
	return scheduler
}

// AddJobInfo add job info with endPoint
func (s *Scheduler) AddJobInfo(endPoint string, jobCount int) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.resources[endPoint]
	if !isExists {
		resource := &ResourceInfo{EndPoint: endPoint, JobCount: jobCount}
		resource.refreshLoadValue()
		s.resources[endPoint] = resource
	} else {
		resource.JobCount += 1
		resource.refreshLoadValue()
		s.resources[endPoint] = resource
	}
}

// RefreshResource refresh resource value
func (s *Scheduler) SetResource(endPoint string, cpuRate int, memoryRate int, jobCount int) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource := &ResourceInfo{EndPoint: endPoint, CpuRate: cpuRate, MemoryRate: memoryRate, JobCount: jobCount}
	resource.refreshLoadValue()
	s.resources[endPoint] = resource
}

// GetResources return scheduler's resources
func (s *Scheduler) Resources() map[string]*ResourceInfo {
	defer s.resourceLocker.RUnlock()
	s.resourceLocker.RLock()
	return s.resources
}

// LoadResource load resource by endPoint
func (s *Scheduler) LoadResource(endPoint string) *ResourceInfo {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.resources[endPoint]
	if !isExists {
		resource := &ResourceInfo{}
		resource.refreshLoadValue()
		s.resources[endPoint] = resource
	}
	return resource
}

// Schedule
func (s *Scheduler) Schedule(balanceType int) ([]*ResourceInfo, error) {
	if s.Resources() == nil || len(s.Resources()) <= 0 {
		return nil, ErrorNotExistsEnoughResource
	}
	if balanceType == Balance_LowerLoad {
		rs := s.getSortLoadResources(s.Resources())
		return rs, nil
	}

	if balanceType == Balance_CpuRate {
		rs := s.getSortCpuResources(s.Resources())
		return rs, nil
	}

	if balanceType == Balance_MemoryRate {
		rs := s.getSortMemoryResources(s.Resources())
		return rs, nil
	}

	if balanceType == Balance_JobCount {
		rs := s.getSortJobResources(s.Resources())
		return rs, nil
	}

	return nil, ErrorNotSupportBalanceType
}

func (s *Scheduler) getSortLoadResources(resources map[string]*ResourceInfo) LoadResources {
	var loadResources LoadResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortCpuResources(resources map[string]*ResourceInfo) CpuResources {
	var loadResources CpuResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortMemoryResources(resources map[string]*ResourceInfo) MemoryResources {
	var loadResources MemoryResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortJobResources(resources map[string]*ResourceInfo) JobResources {
	var loadResources JobResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}
