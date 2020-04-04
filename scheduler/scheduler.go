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
	Scheduler struct {
		resources          map[string]*core.ResourceInfo
		resourceLocker     *sync.RWMutex
		onlineSubmits      map[string]*core.SubmitInfo
		onlineSubmitLocker *sync.RWMutex
	}
)

var ErrorNotSupportBalanceType = errors.New("not support balance type")
var ErrorNotExistsEnoughResource = errors.New("not exists enough resource")

func NewScheduler() *Scheduler {
	scheduler := new(Scheduler)
	scheduler.resources = make(map[string]*core.ResourceInfo)
	scheduler.resourceLocker = new(sync.RWMutex)
	scheduler.onlineSubmits = make(map[string]*core.SubmitInfo)
	scheduler.onlineSubmitLocker = new(sync.RWMutex)

	return scheduler
}

// AddOnlineSubmit add submit info which is online
func (s *Scheduler) AddOnlineSubmit(submit *core.SubmitInfo) {

	s.resourceLocker.Lock()
	endPoint := submit.Worker.EndPoint()
	resource, isExists := s.resources[endPoint]
	if !isExists {
		resource := &core.ResourceInfo{EndPoint: endPoint, TaskCount: 1}
		resource.RefreshLoadValue()
		s.resources[endPoint] = resource
	} else {
		resource.TaskCount += 1
		resource.RefreshLoadValue()
		s.resources[endPoint] = resource
	}
	s.resourceLocker.Unlock()

	s.onlineSubmitLocker.Lock()
	s.onlineSubmits[submit.TaskConfig.TaskID] = submit
	s.onlineSubmitLocker.Unlock()

}

// RefreshResource refresh resource value
func (s *Scheduler) SetResource(resource *core.ResourceInfo) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource.RefreshLoadValue()
	s.resources[resource.EndPoint] = resource
}

// GetResources return scheduler's resources
func (s *Scheduler) Resources() map[string]*core.ResourceInfo {
	defer s.resourceLocker.RUnlock()
	s.resourceLocker.RLock()
	return s.resources
}

// LoadResource load resource by endPoint
func (s *Scheduler) LoadResource(endPoint string) *core.ResourceInfo {
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
func (s *Scheduler) Schedule(balanceType int) ([]*core.ResourceInfo, error) {
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

func (s *Scheduler) getSortLoadResources(resources map[string]*core.ResourceInfo) core.LoadResources {
	var loadResources core.LoadResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortCpuResources(resources map[string]*core.ResourceInfo) core.CpuResources {
	var loadResources core.CpuResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortMemoryResources(resources map[string]*core.ResourceInfo) core.MemoryResources {
	var loadResources core.MemoryResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}

func (s *Scheduler) getSortJobResources(resources map[string]*core.ResourceInfo) core.JobResources {
	var loadResources core.JobResources
	for _, resource := range resources {
		loadResources = append(loadResources, resource)
	}
	sort.Sort(loadResources)
	return loadResources
}
