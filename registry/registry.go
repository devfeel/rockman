package registry

import (
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/registry/consul"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	ServerUrl string
	*consul.ConsulClient
}

func NewRegistry(regServer string) (*Registry, error) {
	reg := &Registry{}
	reg.ServerUrl = regServer
	regClient, err := consul.NewConsulClient(regServer)
	if err != nil {
		logger.Node().Debug("Registry init error: " + err.Error())
		logger.Node().Error(err, "Registry init error")
		return nil, err
	}
	logger.Node().Debug("Registry init success.")
	reg.ConsulClient = regClient
	return reg, nil
}

// CreateLocker create locker to registry with key/value
func (r *Registry) CreateLocker(key string, value string, ttl string) (*consul.Locker, error) {
	if ttl == "" {
		ttl = "10s"
	}
	opts := &api.LockOptions{
		Key:   key,
		Value: []byte(value),
		SessionOpts: &api.SessionEntry{
			Name:     key,
			TTL:      ttl,
			Behavior: "delete",
		},
	}
	locker, err := r.CreateLockerOpts(opts)
	if err != nil {
		return nil, err
	}
	return locker, nil
}
