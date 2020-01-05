package webui

import "github.com/devfeel/rockman/src/logger"

type WebServer struct {
}

func NewWebServer() *WebServer {
	s := &WebServer{}
	logger.Default().Debug("WebUI Init Success!")
	return s
}

func (s *WebServer) Start() error {
	return nil
}
