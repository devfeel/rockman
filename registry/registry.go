package registry

import (
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/registry/consul"
	"github.com/devfeel/rockman/util/netx"
	"github.com/hashicorp/consul/api"
	"strconv"
	"time"
)

const defaultRetryPingLimit = 3

type (
	Registry struct {
		ServerUrl string
		*consul.ConsulClient
		OnServerOnline  NetChangeHandle
		OnServerOffline NetChangeHandle
		isServerOnline  bool
		isStart         bool
	}

	NetChangeHandle func()
)

func NewRegistry(regServer string) (*Registry, error) {
	reg := &Registry{}
	reg.ServerUrl = regServer
	reg.isServerOnline = true
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

func (r *Registry) Start() error {
	if r.isStart {
		return nil
	}
	logger.Default().Debug("Registry start...")
	r.isStart = true
	r.watchPingRegistry()
	return nil
}

func (r *Registry) Stop() error {
	logger.Default().Debug("Registry Stop...")
	r.isStart = false
	return nil
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

// watchPingRegistry
// check connect registry server
func (r *Registry) watchPingRegistry() {
	lt := "Registry.watchPingRegistry "
	logger.Default().Debug(lt + "running...")

	doQuery := func() bool {
		// check connect registry server
		result := netx.CheckTcpConnect(r.ServerUrl)
		return result
	}

	go func() {
		var retryCount int
		limit := defaultRetryPingLimit
		for {
			if !r.isStart {
				logger.Default().DebugS(lt + "registry is not start, now stop watch.")
				return
			}
			time.Sleep(time.Second * 10)
			result := doQuery()
			if !result {
				if !r.isServerOnline {
					logger.Default().DebugS(lt + "ping registry failed.")
					continue
				} else {
					logger.Default().DebugS(lt + "ping registry failed[" + strconv.Itoa(retryCount) + "].")
					retryCount += 1
				}
				if retryCount > limit {
					logger.Default().DebugS(lt + "retry count more than " + strconv.Itoa(limit) + ", now is confirm unable to connect registry.")
					r.isServerOnline = false
					if r.OnServerOffline != nil {
						r.OnServerOffline()
					}
				}
			} else {
				if !r.isServerOnline {
					logger.Default().DebugS(lt + "ping registry success, now change server online.")
					r.isServerOnline = true
					if r.OnServerOnline != nil {
						r.OnServerOnline()
					}
				}
				retryCount = 0
			}
		}
	}()
}
