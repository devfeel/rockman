package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman-webui/src/protected/viewModel/user"
)

type UserController struct {
}

func (c *UserController) Login(ctx dotweb.Context) error {
	userName := ctx.QueryString("UserName")
	userPwd := ctx.QueryString("UserPwd")
	loginUser := user.LoginUserDto{}
	loginUser.Token = userName + "|" + userPwd
	loginUser.UserName = userName
	if userName == "root" && userPwd == "root" {
		return ctx.WriteJson(SuccessResponse(loginUser))
	}
	return ctx.WriteJson(FailedResponse(-1000, "password error!"))
}
