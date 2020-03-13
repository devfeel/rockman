package rpc

import (
	"github.com/devfeel/rockman/src/core/packets"
	"testing"
)

const (
	serverUrl = "127.0.0.1:2398"
)

func TestRpcClient_CallEcho(t *testing.T) {
	client := getRpcClient()
	message := "echo message"
	wantMessage := "echo message"
	err, result := client.CallEcho(message)
	if err != nil {
		t.Error(err)
	} else {
		if result == wantMessage {
			t.Log("success:", wantMessage, result)
		} else {
			t.Error("failed:", wantMessage, result)
		}
	}
}

func TestRpcClient_CallRegisterNode(t *testing.T) {
	client := getRpcClient()
	nodeInfo := packets.NodeInfo{Host: "127.0.0.1", Port: "2401", NodeID: "TestNode", IsMaster: true, IsWorker: true}
	err, result := client.CallRegisterNode(nodeInfo)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func getRpcClient() *RpcClient {
	return NewRpcClient(serverUrl)
}
