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
	executorController = new(controllers.ExecutorController)
	nodeController     = new(controllers.NodeController)
	clusterController  = new(controllers.ClusterController)
	userController     = new(controllers.UserController)
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
	g := s.webApp.HttpServer.Group("/task")
	g.GET("/list", executorController.ShowExecutors)
	g.GET("/logs", executorController.ShowExecLogs)

	g = s.webApp.HttpServer.Group("/node")
	g.GET("/list", nodeController.ShowNodes)

	g = s.webApp.HttpServer.Group("/cluster")
	g.GET("/resources", clusterController.ShowResources)
	g.GET("/executors", clusterController.ShowExecutors)

	g = s.webApp.HttpServer.Group("/user")
	g.GET("/login", userController.Login)

	s.webApp.HttpServer.Router().ServerFile("/*filepath", "/wwwroot/")

}
