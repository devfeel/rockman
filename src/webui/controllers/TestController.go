package controllers

import "github.com/devfeel/dotweb"

type TestController struct {
}

func (c *TestController) Echo(ctx dotweb.Context) error {
	return ctx.WriteString(ctx.Request().QueryString("message"))
}
