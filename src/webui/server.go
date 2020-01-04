package webui

type WebServer struct {
}

func NewWebServer() *WebServer {
	return &WebServer{}
}

func (s *WebServer) Start() error {
	return nil
}
