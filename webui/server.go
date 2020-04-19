package webui

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/cors"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/node"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/controllers"
	"strings"
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
	s.initModule()
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

	executorController := controllers.NewExecutorController()
	nodeController := new(controllers.NodeController)
	clusterController := new(controllers.ClusterController)
	userController := new(controllers.UserController)

	g := s.webApp.HttpServer.Group("/api/task")
	g.GET("/list", executorController.ShowExecutors)
	g.POST("/save", executorController.SaveExecutor)
	g.POST("/update", executorController.UpdateExecutor)
	g.GET("/get", executorController.QueryById)
	g.POST("/execlogs", executorController.ShowExecLogs)
	g.POST("/statelogs", executorController.ShowStateLog)

	g = s.webApp.HttpServer.Group("/api/node")
	g.GET("/list", nodeController.ShowNodes)

	g = s.webApp.HttpServer.Group("/api/cluster")
	g.GET("/resources", clusterController.ShowResources)
	g.GET("/executors", clusterController.ShowExecutors)
	g.GET("/info", clusterController.ShowClusterInfo)

	g = s.webApp.HttpServer.Group("/api/user")
	g.GET("/login", userController.Login)

	s.webApp.HttpServer.ServerFile("/static/*", "./webapp/")

}

func (s *WebServer) initModule() {
	s.webApp.HttpServer.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			path := ctx.Request().URL.Path
			if strings.HasPrefix(path, "/static/") && !strings.HasPrefix(path, "/static/static") {
				ctx.Request().Request.URL.Path = "/static" + path
			}
		},
	})
}
