package scheduler

import (
	"errors"
	"github.com/devfeel/rockman/core"
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
	Scheduler interface {
		SetResource(resource *core.ResourceInfo)
		Resources() map[string]*core.ResourceInfo
		LoadResource(endPoint string) *core.ResourceInfo
		Schedule(balanceType int) ([]*core.ResourceInfo, error)
	}

	StandardScheduler struct {
		resources      map[string]*core.ResourceInfo
		resourceLocker *sync.RWMutex
	}
)

var ErrorNotSupportBalanceType = errors.New("not support balance type")
var ErrorNotExistsEnoughResource = errors.New("not exists enough resource")

func NewScheduler() Scheduler {
	scheduler := new(StandardScheduler)
	scheduler.resources = make(map[string]*core.ResourceInfo)
	scheduler.resourceLocker = new(sync.RWMutex)

	return scheduler
}

// RefreshResource refresh resource value
func (s *StandardScheduler) SetResource(resource *core.ResourceInfo) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource.RefreshLoadValue()
	s.resources[resource.EndPoint] = resource
}

// GetResources return scheduler's resources
func (s *StandardScheduler) Resources() map[string]*core.ResourceInfo {
	defer s.resourceLocker.RUnlock()
	s.resourceLocker.RLock()
	return s.resources
}

// LoadResource load resource by endPoint
func (s *StandardScheduler) LoadResource(endPoint string) *core.ResourceInfo {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.resources[endPoint]
	if !isExists {
		resource := &core.ResourceInfo{}
		resource.RefreshLoadValue()
		s.resources[endPoint] = resource
	}
	return resource
}

// Schedule
func (s *StandardScheduler) Schedule(balanceType int) ([]*core.ResourceInfo, error) {
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

func (s *StandardScheduler) getSortLoadResources(resources map[string]*core.ResourceInfo) core.LoadResources {
	var loadResources core.LoadResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *StandardScheduler) getSortCpuResources(resources map[string]*core.ResourceInfo) core.CpuResources {
	var loadResources core.CpuResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *StandardScheduler) getSortMemoryResources(resources map[string]*core.ResourceInfo) core.MemoryResources {
	var loadResources core.MemoryResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *StandardScheduler) getSortJobResources(resources map[string]*core.ResourceInfo) core.JobResources {
	var loadResources core.JobResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}
