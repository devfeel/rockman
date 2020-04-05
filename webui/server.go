package webui

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/cors"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/controllers"
)

var (
	testController = new(controllers.TestController)
	taskController = new(controllers.TaskController)
	nodeController = new(controllers.NodeController)
)

type WebServer struct {
	webApp     *dotweb.DotWeb
	listenAddr string
}

func NewWebServer(logPath string, node *node.Node) *WebServer {
	s := &WebServer{}
	s.webApp = dotweb.New()
	s.webApp.SetLogPath(logPath)
	s.webApp.SetEnabledLog(true)
	s.webApp.UseRequestLog()
	s.webApp.Use(cors.Middleware(cors.NewConfig().UseDefault()))
	s.webApp.Items.Set(_const.ItemKey_Node, node)
	s.initRoute()
	logger.Default().Debug("WebUI init success.")
	return s
}

func (s *WebServer) ListenAndServe(listenAddr string) error {
	s.listenAddr = listenAddr
	logger.Default().Debug("WebServer begin listen " + s.listenAddr)
	err := s.webApp.ListenAndServe(s.listenAddr)
	if err != nil {
		return err
	}
	return nil
}

func (s *WebServer) initRoute() {
	g := s.webApp.HttpServer.Group("/test")
	g.GET("/index", testController.Echo)

	g = s.webApp.HttpServer.Group("/task")
	g.GET("/list", taskController.ShowTasks)
	g.GET("/logs", taskController.ShowLogs)

	g = s.webApp.HttpServer.Group("/node")
	g.GET("/list", nodeController.ShowNodeList)
	g.GET("/resource", nodeController.ShowResource)
}
