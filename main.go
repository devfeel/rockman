package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/rpc"
	"github.com/devfeel/rockman/util/exception"
	"github.com/devfeel/rockman/webui"
	"os"
	"time"
)

var CurNode *node.Node
var CurRpcServer *rpc.RpcServer
var CurWebServer *webui.WebServer

const (
	ProjectName = "rockman"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			exMsg := exception.CatchError(ProjectName+":main", err)
			logger.Default().Error(errors.New(exMsg), "")
			os.Stdout.Write([]byte(exMsg))
		}
	}()

	printLogo()
	var err error

	// load config file
	profile := config.SingleNodeProfile()

	// start log service
	logger.StartLogService("config")

	//start worker node
	CurNode, err = node.NewNode(profile)
	if err != nil {
		fmt.Println(err)
		return
	}

	//start rpc server
	CurRpcServer = rpc.NewRpcServer(profile, CurNode)
	//start web server
	if profile.Node.IsMaster {
		CurWebServer = webui.NewWebServer(profile.Logger.LogPath, CurNode)
	}

	go func() {
		err := CurRpcServer.Listen()
		if err != nil {
			logger.Default().Error(err, "RpcServer start error")
			panic(errors.New("RpcServer start error"))
		}
	}()

	if profile.Node.IsMaster {
		go func() {
			err := CurWebServer.ListenAndServe(profile.WebUI.HttpHost + ":" + profile.WebUI.HttpPort)
			if err != nil {
				logger.Default().Error(err, "WebUI start error")
				panic(errors.New("WebUI start error"))
			}
		}()
	}

	time.Sleep(time.Second)
	//start node
	CurNode.Start()

	for {
		time.Sleep(time.Hour)
	}
}

func printLogo() {
	fmt.Println(".______      ______     ______  __  ___ .___  ___.      ___      .__   __. ")
	fmt.Println("|   _  \\    /  __  \\   /      ||  |/  / |   \\/   |     /   \\     |  \\ |  | ")
	fmt.Println("|  |_)  |  |  |  |  | |  ,----'|  '  /  |  \\  /  |    /  ^  \\    |   \\|  | ")
	fmt.Println("|      /   |  |  |  | |  |     |    <   |  |\\/|  |   /  /_\\  \\   |  . `  | ")
	fmt.Println("|  |\\  \\--.|  `--'  | |  `----.|  .  \\  |  |  |  |  /  _____  \\  |  |\\   | ")
	fmt.Println("| _| `.___| \\______/   \\______||__|\\__\\ |__|  |__| /__/     \\__\\ |__| \\__| ")
}
