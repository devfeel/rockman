package client

import (
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/runtime/executor"
	"testing"
)

const (
	serverUrl = "116.62.16.66:2398"
	//serverUrl = "118.31.32.168:2398"
	//serverUrl = "127.0.0.1:2398"
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
	worker := &core.NodeInfo{Host: "127.0.0.1", Port: "2401", NodeID: "TestNode"}
	err, result := client.CallRegisterNode(worker)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallQueryNodes(t *testing.T) {
	client := getRpcClient()
	page := &core.PageInfo{PageIndex: 1, PageSize: 10}
	err, result := client.CallQueryNodes(page)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRegisterHttpExecutor(t *testing.T) {
	client := getRpcClient()
	conf := &core.TaskConfig{}
	conf.TaskID = "TestRpcClient-http-debug"
	conf.TaskType = "cron"
	conf.TargetType = "http"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"
	conf.TargetConfig = &executor.HttpConfig{
		Url:    "http://www.baidu.com",
		Method: "HEAD",
	}

	err, result := client.CallRegisterExecutor(conf)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRegisterShellScriptExecutor(t *testing.T) {
	client := getRpcClient()
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-shell-Script"
	conf.TaskType = "cron"
	conf.TargetType = "shell"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.ShellConfig{
		Script: "echo ok",
		Type:   "Script",
	}

	err, result := client.CallRegisterExecutor(conf)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRegisterShellFileExecutor(t *testing.T) {
	client := getRpcClient()
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-shell-File"
	conf.TaskType = "cron"
	conf.TargetType = "shell"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.ShellConfig{
		Script: "hello.sh",
		Type:   "File",
	}

	err, result := client.CallRegisterExecutor(conf)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallRegisterGoExecutor(t *testing.T) {
	client := getRpcClient()
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-shell-Go"
	conf.TaskType = "cron"
	conf.TargetType = "goso"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.GoConfig{
		FileName: "plugin.so",
	}

	err, result := client.CallRegisterExecutor(conf)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallSubmitHttpExecutor(t *testing.T) {
	client := getRpcClient()
	submit := new(core.SubmitInfo)
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-http"
	conf.TaskType = "cron"
	conf.TargetType = "http"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = "http-url"
	conf.TargetConfig = &executor.HttpConfig{
		Url:    "http://www.baidu.com",
		Method: "HEAD",
	}

	submit.TaskConfig = conf
	submit.Worker = &core.NodeInfo{
		Host: "118.31.32.168",
		Port: "2398",
	}

	err, result := client.CallSubmitExecutor(submit)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallSubmitShellScriptExecutor(t *testing.T) {
	client := getRpcClient()
	submit := new(core.SubmitInfo)
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-shell-Script"
	conf.TaskType = "cron"
	conf.TargetType = "shell"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.ShellConfig{
		Script: "echo ok",
		Type:   "Script",
	}

	submit.TaskConfig = conf
	submit.Worker = &core.NodeInfo{
		Host: "118.31.32.168",
		Port: "2398",
	}

	err, result := client.CallSubmitExecutor(submit)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallSubmitShellFileExecutor(t *testing.T) {
	client := getRpcClient()
	submit := new(core.SubmitInfo)
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-shell-File"
	conf.TaskType = "cron"
	conf.TargetType = "shell"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.ShellConfig{
		Script: "hello.sh",
		Type:   "File",
	}

	submit.TaskConfig = conf
	submit.Worker = &core.NodeInfo{
		Host: "118.31.32.168",
		Port: "2398",
	}

	err, result := client.CallSubmitExecutor(submit)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallSubmitGoExecutor(t *testing.T) {
	client := getRpcClient()
	submit := new(core.SubmitInfo)
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-GoSo"
	conf.TaskType = "cron"
	conf.TargetType = "goso"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.GoConfig{
		FileName: "plugin.so",
	}

	submit.TaskConfig = conf
	submit.Worker = &core.NodeInfo{
		Host: "118.31.32.168",
		Port: "2398",
	}

	err, result := client.CallSubmitExecutor(submit)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success:", result)
	}
}

func TestRpcClient_CallSubmitLBGoExecutor(t *testing.T) {
	client := getRpcClient()
	submit := new(core.SubmitInfo)
	conf := &core.TaskConfig{}
	conf.TaskID = "Test-GoSo-LB"
	conf.TaskType = "cron"
	conf.TargetType = "goso"
	conf.IsRun = true
	conf.DueTime = 0
	conf.Interval = 0
	conf.Express = "0 * * * * *"
	conf.TaskData = ""
	conf.TargetConfig = &executor.GoConfig{
		FileName: "plugin.so",
	}
	submit.TaskConfig = conf
	err, result := client.CallSubmitExecutor(submit)
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
	taskId := "Test-shell-File"

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
	return NewRpcClient(serverUrl, false, "", "")
}
