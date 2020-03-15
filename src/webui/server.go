package webui

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/cors"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/webui/controllers"
)

var (
	testController = new(controllers.TestController)
	taskController = new(controllers.TaskController)
)

type WebServer struct {
	webApp     *dotweb.DotWeb
	listenAddr string
}

func NewWebServer(logPath string) *WebServer {
	s := &WebServer{}
	s.webApp = dotweb.New()
	s.webApp.SetLogPath(logPath)
	s.webApp.SetEnabledLog(true)
	s.webApp.UseRequestLog()
	s.webApp.Use(cors.Middleware(cors.NewConfig().UseDefault()))
	s.initRoute()
	logger.Default().Debug("WebUI init success.")
	return s
}

func (s *WebServer) ListenAndServe(listenAddr string) error {
	s.listenAddr = listenAddr
	logger.Default().Debug("WebServer.StartServer => " + s.listenAddr)
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
	g.GET("/taskbynode", taskController.ShowTaskListByNodeID)
}
