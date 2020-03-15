package node

import (
	"fmt"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/util/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"time"
)

type Registry struct {
	ServerUrl             string
	RegServer             *consul.ConsulClient
	LeaderKey             string
	LeaderServer          string
	lastGetLeaderInfoTime time.Time
}

// initRegistry init Registry and reg server
func initRegistry(registerServer string, leaderKey string) (*Registry, error) {
	registry := new(Registry)
	registry.ServerUrl = registerServer
	registry.LeaderKey = leaderKey
	regServer, err := consul.NewConsulClient(registerServer)
	if err != nil {
		logger.Node().Debug(fmt.Sprint("Registry init error", err.Error()))
		logger.Node().Error(err, "Registry init error")
		return nil, err
	}
	registry.RegServer = regServer
	logger.Node().Debug("Registry init success.")
	return registry, nil
}

// electionLeader election leader role to registry server
func (r *Registry) electionLeader(leaderServer string, checkUrl string) (bool, error) {
	opts := &api.LockOptions{
		Key:         r.LeaderKey,
		Value:       []byte(leaderServer),
		SessionTTL:  "10s",
		SessionName: leaderServer,
	}
	locker, err := r.RegServer.CreateLockerOpts(opts)
	if err != nil {
		return false, err
	}

	_, err = locker.Locker.Lock(nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLeaderInfo get leader info from leader key
// must check is locked by leader session
// cache in memory with 1 minute
func (r *Registry) GetLeaderInfo() (string, error) {
	if r.LeaderServer != "" && time.Now().Sub(r.lastGetLeaderInfoTime) < time.Minute {
		return r.LeaderServer, nil
	}
	kvPair, err := r.RegServer.Get(r.LeaderKey)
	if err != nil {
		return "", err
	} else {
		if kvPair.Session == "" {
			return "", errors.New("no leader info exists")
		} else {
			r.LeaderServer = string(kvPair.Value)
			r.lastGetLeaderInfoTime = time.Now()
			return r.LeaderServer, nil
		}
	}
}
