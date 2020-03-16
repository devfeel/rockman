package client

import (
	"github.com/devfeel/rockman/src/packets"
	"github.com/devfeel/rockman/src/runtime/executor"
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

func TestRpcClient_CallRegisterWorker(t *testing.T) {
	client := getRpcClient()
	worker := &packets.WorkerInfo{Host: "127.0.0.1", Port: "2401", NodeID: "TestNode"}
	err, result := client.CallRegisterWorker(worker)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRegisterExecutor(t *testing.T) {
	client := getRpcClient()
	conf := &executor.HttpTaskConfig{}
	conf.TaskID = "TestRpcClient-http-debug"
	conf.TaskType = "cron"
	conf.TargetType = "http"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"

	err, result := client.CallRegisterExecutor(conf)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func getRpcClient() *RpcClient {
	return NewRpcClient(serverUrl)
}
