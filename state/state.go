package state

import (
	"sync"
)

type (
	State struct {
		Resources      map[string]*ResourceInfo
		resourceLocker *sync.RWMutex
	}

	ResourceInfo struct {
		EndPoint   string
		CpuRate    int //cpu rate, refresh per 1 minute
		MemoryRate int //memory rate, refresh per 1 minute
		JobCount   int //job count
		LoadValue  int //load value = cpu * 30 + memory * 30 + jobs * 40
	}
)

func NewState() *State {
	state := new(State)
	state.Resources = make(map[string]*ResourceInfo)
	state.resourceLocker = new(sync.RWMutex)
	return state
}

// AddJob add job to endPoint
func (s *State) AddJob(endPoint string, jobCount int) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.Resources[endPoint]
	if !isExists {
		resource := &ResourceInfo{EndPoint: endPoint, JobCount: jobCount}
		resource.refreshLoadValue()
		s.Resources[endPoint] = resource
	} else {
		resource.JobCount += 1
		resource.refreshLoadValue()
		s.Resources[endPoint] = resource
	}
}

// RefreshResource refresh resource value
func (s *State) RefreshResource(endPoint string, cpuRate int, memoryRate int, jobCount int) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource := &ResourceInfo{EndPoint: endPoint, CpuRate: cpuRate, MemoryRate: memoryRate, JobCount: jobCount}
	resource.refreshLoadValue()
	s.Resources[endPoint] = resource
}

// LoadResource load resource by endPoint
func (s *State) LoadResource(endPoint string) *ResourceInfo {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.Resources[endPoint]
	if !isExists {
		resource := &ResourceInfo{}
		resource.refreshLoadValue()
		s.Resources[endPoint] = resource
	}
	return resource
}

// refreshLoadValue refresh resource's load value
func (r *ResourceInfo) refreshLoadValue() int {
	r.LoadValue = r.CpuRate*30 + r.MemoryRate*30 + r.JobCount*40
	return r.LoadValue
}
