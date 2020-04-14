package node

import (
	"github.com/devfeel/rockman/core"
	"time"
)

func (n *Node) getExecutorInitFlag() (bool, error) {
	kv, _, err := n.Registry.Get(getExecutorInitFlagKey(n.ClusterId()), nil)
	if err == nil {
		return false, err
	}
	if kv == nil {
		return false, nil
	}
	return true, nil
}

func (n *Node) deleteExecutorInitFlag() error {
	_, err := n.Registry.Delete(getExecutorInitFlagKey(n.ClusterId()), nil)
	return err
}

func (n *Node) setExecutorInitFlag() error {
	_, err := n.Registry.Set(getExecutorInitFlagKey(n.ClusterId()), "true", nil)
	return err
}

func (n *Node) setExecutorChangeFlag() error {
	_, err := n.Registry.Set(getExecutorChangeFlagKey(n.ClusterId()), time.Now().Format(core.DefaultTimeLayout), nil)
	return err
}

func getExecutorInitFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/executor-init"
}

func getExecutorChangeFlagKey(clusterId string) string {
	return core.ClusterKeyPrefix + clusterId + "/flags/executor-change"
}
