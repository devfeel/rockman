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
	s.webApp.Items.Set(_const.ItemKeyNode, node)

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

	g := s.webApp.HttpServer.Group("/api/task")
	g.POST("/list", executorController.ShowExecutors)
	g.POST("/save", executorController.SaveExecutor)
	g.POST("/update", executorController.UpdateExecutor)
	g.GET("/get", executorController.QueryById)
	g.GET("/delete", executorController.DeleteById)
	g.POST("/execlogs", executorController.ShowExecLogs)
	g.POST("/statelogs", executorController.ShowStateLog)

	g = s.webApp.HttpServer.Group("/api/node")
	g.POST("/list", nodeController.ShowNodes)

	g = s.webApp.HttpServer.Group("/api/cluster")
	g.GET("/resources", clusterController.ShowResources)
	g.GET("/executors", clusterController.ShowExecutors)
	g.GET("/info", clusterController.ShowClusterInfo)

	g = s.webApp.HttpServer.Group("/api/user")
	g.GET("/login", userController.Login)

	// g = s.webApp.HttpServer.Group("/*")

	s.webApp.HttpServer.Router().ServerFile("/static/*filepath", "release/wwwroot/")

}
