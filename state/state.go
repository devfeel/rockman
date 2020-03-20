package state

import "sync"

type (
	State struct {
		Resources      map[string]*ResourceInfo
		resourceLocker *sync.RWMutex
	}

	ResourceInfo struct {
		CpuRate    int //cpu rate, refresh per 1 minute
		MemoryRate int //memory rate, refresh per 1 minute
		JobCount   int //job count
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
		s.Resources[endPoint] = &ResourceInfo{JobCount: jobCount}
	} else {
		resource.JobCount += 1
		s.Resources[endPoint] = resource
	}
}

// RefreshResource refresh resource value
func (s *State) RefreshResource(endPoint string, cpuRate int, memoryRate int, jobCount int) {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource := &ResourceInfo{CpuRate: cpuRate, MemoryRate: memoryRate, JobCount: jobCount}
	s.Resources[endPoint] = resource
}

// LoadResource load resource by endPoint
func (s *State) LoadResource(endPoint string) *ResourceInfo {
	defer s.resourceLocker.Unlock()
	s.resourceLocker.Lock()
	resource, isExists := s.Resources[endPoint]
	if !isExists {
		resource := &ResourceInfo{}
		s.Resources[endPoint] = resource
	}
	return resource
}
