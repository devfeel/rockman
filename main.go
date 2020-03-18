package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/rpc"
	"github.com/devfeel/rockman/webui"
	"time"
)

var CurNode *node.Node
var CurRpcServer *rpc.RpcServer
var CurWebServer *webui.WebServer

func main() {
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
		CurWebServer = webui.NewWebServer(profile.Logger.LogPath)
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

	//start node
	CurNode.Start()

	for {
		time.Sleep(time.Hour)
	}
}

func printLogo() {
	fmt.Println(".______        ______     ______  __  ___ .___  ___.      ___      .__   __. ")
	fmt.Println("|   _  \\      /  __  \\   /      ||  |/  / |   \\/   |     /   \\     |  \\ |  | ")
	fmt.Println("|  |_)  |    |  |  |  | |  ,----'|  '  /  |  \\  /  |    /  ^  \\    |   \\|  | ")
	fmt.Println("|      /     |  |  |  | |  |     |    <   |  |\\/|  |   /  /_\\  \\   |  . `  | ")
	fmt.Println("|  |\\  \\----.|  `--'  | |  `----.|  .  \\  |  |  |  |  /  _____  \\  |  |\\   | ")
	fmt.Println("| _| `._____| \\______/   \\______||__|\\__\\ |__|  |__| /__/     \\__\\ |__| \\__| ")
}
