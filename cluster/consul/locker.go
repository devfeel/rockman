package consul

import "github.com/hashicorp/consul/api"

type Locker struct {
	Locker *api.Lock
}

func (l *Locker) Lock() (<-chan struct{}, error) {
	return l.Locker.Lock(nil)
}

func (l *Locker) UnLock() error {
	return l.Locker.Unlock()
}
