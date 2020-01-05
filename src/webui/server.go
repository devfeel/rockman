package webui

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/middleware/cors"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/logger"
	"github.com/devfeel/rockman/src/webui/controllers"
	"strconv"
)

var (
	testController = new(controllers.TestController)
)

type WebServer struct {
	webApp     *dotweb.DotWeb
	listenAddr string
}

func NewWebServer() *WebServer {
	s := &WebServer{}
	s.listenAddr = config.CurrentProfile.Node.HttpHost + ":" + strconv.Itoa(config.CurrentProfile.Node.HttpPort)
	s.webApp = dotweb.New()
	s.webApp.SetLogPath(config.CurrentProfile.Logger.LogPath + "/webui/")
	s.webApp.SetEnabledLog(true)
	s.webApp.UseRequestLog()
	s.webApp.Use(cors.Middleware(cors.NewConfig().UseDefault()))
	s.initRoute()
	logger.Default().Debug("WebUI Init Success!")
	return s
}

func (s *WebServer) Start() error {
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
}
