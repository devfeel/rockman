package client

import (
	"github.com/devfeel/rockman/packets"
	"github.com/devfeel/rockman/runtime/executor"
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
	worker := &packets.NodeInfo{Host: "127.0.0.1", Port: "2401", NodeID: "TestNode"}
	err, result := client.CallRegisterNode(worker)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallQueryNodes(t *testing.T) {
	client := getRpcClient()
	page := &packets.PageInfo{PageIndex: 1, PageSize: 10}
	err, result := client.CallQueryNodes(page)
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

func TestRpcClient_CallStartExecutor(t *testing.T) {
	client := getRpcClient()
	taskId := "TestRpcClient-http-debug"

	err, result := client.CallStartExecutor(taskId)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallStopExecutor(t *testing.T) {
	client := getRpcClient()
	taskId := "TestRpcClient-http-debug"

	err, result := client.CallStopExecutor(taskId)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRemoveExecutor(t *testing.T) {
	client := getRpcClient()
	taskId := "TestRpcClient-http-debug"

	err, result := client.CallRemoveExecutor(taskId)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallQueryExecutors(t *testing.T) {
	client := getRpcClient()

	err, result := client.CallQueryExecutorConfig("")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func getRpcClient() *RpcClient {
	return NewRpcClient(serverUrl)
}
