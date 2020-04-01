package main

import (
	"errors"
	"flag"
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
	ProjectName  = "rockman"
	cmdNodeType  = "nodetype"
	cmdOuterHost = "outerhost"
	cmdOuterPort = "outerport"
	cmdCluster   = "cluster"

	version = "2020.401"
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
	profile := config.DefaultProfile()

	parseFlag(profile)

	// start log service
	logger.StartLogService("config")

	shutdownChan := make(chan string)

	//start worker node
	CurNode, err = node.NewNode(profile, shutdownChan)
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
	err = CurNode.Start()
	if err != nil {
		logger.Default().Error(err, "Node start error")
		return
	}
	logger.Default().Debug("Node start success, service running...")

	for {
		<-shutdownChan
		logger.Default().Debug("Node Shutdown.")
	}
	logger.Default().Debug("Node Close.")
}

func parseFlag(profile *config.Profile) {
	var nodeType, outerHost, outerPort, cluster string
	flag.StringVar(&nodeType, cmdNodeType, "", "node type, full or master or worker")
	flag.StringVar(&outerHost, cmdOuterHost, "", "node outer host")
	flag.StringVar(&outerPort, cmdOuterPort, "", "node outer port")
	flag.StringVar(&cluster, cmdCluster, "", "node cluster id")

	flag.Parse()
	if nodeType == "master" {
		profile.Node.IsWorker = false
	}

	if nodeType == "worker" {
		profile.Node.IsMaster = false
	}

	profile.Rpc.OuterHost = outerHost
	profile.Rpc.OuterPort = outerPort

	if cluster != "" {
		profile.Cluster.ClusterId = cluster
	}
}

func printLogo() {
	fmt.Println(".______      ______     ______  __  ___ .___  ___.      ___      .__   __. ")
	fmt.Println("|   _  \\    /  __  \\   /      ||  |/  / |   \\/   |     /   \\     |  \\ |  | ")
	fmt.Println("|  |_)  |  |  |  |  | |  ,----'|  '  /  |  \\  /  |    /  ^  \\    |   \\|  | ")
	fmt.Println("|      /   |  |  |  | |  |     |    <   |  |\\/|  |   /  /_\\  \\   |  . `  | ")
	fmt.Println("|  |\\  \\--.|  `--'  | |  `----.|  .  \\  |  |  |  |  /  _____  \\  |  |\\   | ")
	fmt.Println("| _| `.___| \\______/   \\______||__|\\__\\ |__|  |__| /__/     \\__\\ |__| \\__| ")
	fmt.Printf("%c[1m%s%c[0m\n", 0x1B, "                                                       Version: "+version, 0x1B)
}
