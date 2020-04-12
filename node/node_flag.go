package node

import (
	"github.com/devfeel/rockman/core"
	"time"
)

func (n *Node) getInitFlag() (bool, error) {
	kv, _, err := n.Registry.Get(getInitFlagKey(n.ClusterId()), nil)
	if err == nil {
		return false, err
	}
	if kv == nil {
		return false, nil
	}
	return true, nil
}

func (n *Node) setInitFlag() error {
	_, err := n.Registry.Set(getInitFlagKey(n.ClusterId()), "true", nil)
	return err
}

func (n *Node) setExecutorChangeFlag() error {
	_, err := n.Registry.Set(getExecutorChangeFlagKey(n.ClusterId()), time.Now().Format(core.DefaultTimeLayout), nil)
	return err
}

func getInitFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/init"
}

func getExecutorChangeFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/executor-change"
}
