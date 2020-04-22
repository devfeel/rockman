package middleware

import (
	"fmt"

	"github.com/devfeel/dotweb"
)

type TestMiddleware struct {
	dotweb.BaseMiddlware
}

func (m *TestMiddleware) Handle(ctx dotweb.Context) error {
	fmt.Println(1)
	err := m.Next(ctx)
	response := ctx.Response()
	fmt.Println(response.Status)
	if response.Status == 404 {
		ctx.Request().Request.URL.Path = "/static/index.html"
		ctx.WriteStringC(200, nil)
		ctx.Response().SetStatusCode(200)
	}
	return err
}

// Middleware create new CORS Middleware
func Middleware() *TestMiddleware {
	return &TestMiddleware{}
}

// DefaultMiddleware create new CORS Middleware with default config
func DefaultMiddleware() *TestMiddleware {
	return &TestMiddleware{}
}
