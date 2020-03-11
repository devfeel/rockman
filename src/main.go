package main

import (
	"fmt"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/node"
	"time"
)

var CurrentNode *node.Node

func main() {
	println("Welcome to RockMan!")

	var err error

	// load config file
	profile := config.SingleNodeProfile()

	// start log service
	logger.StartLogService("config")

	//start worker node
	CurrentNode, err = node.NewNode(profile)
	if err != nil {
		fmt.Println(err)
		return
	}
	CurrentNode.Start()

	for {
		time.Sleep(time.Hour)
	}
}
