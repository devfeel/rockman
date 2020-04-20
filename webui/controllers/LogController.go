package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/protected/service"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

type LogController struct {
}

func NewLogController() *LogController {
	return &LogController{}
}

// ShowExecLogs
func (c *LogController) ShowTaskExecLogs(ctx dotweb.Context) error {
	qr := new(contract.TaskExecLogQR)
	err := ctx.Bind(qr)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}
	if qr.PageSize <= 0 {
		qr.PageSize = _const.DefaultPageSize
	}
	logService := service.NewLogService()
	result, err := logService.QueryExecLogs(qr.TaskID, &qr.PageRequest)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowStateLog
func (c *LogController) ShowTaskStateLog(ctx dotweb.Context) error {
	qr := new(contract.TaskStateLogQR)
	err := ctx.Bind(qr)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}

	if qr.PageSize <= 0 {
		qr.PageSize = _const.DefaultPageSize
	}

	logService := service.NewLogService()
	result, err := logService.QueryStateLog(qr.TaskID, &qr.PageRequest)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowStateLog
func (c *LogController) ShowTaskSubmitLog(ctx dotweb.Context) error {
	qr := new(contract.TaskSubmitLogQR)
	err := ctx.Bind(qr)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}

	if qr.PageSize <= 0 {
		qr.PageSize = _const.DefaultPageSize
	}

	logService := service.NewLogService()
	result, err := logService.QueryTaskSubmitLog(qr.TaskID, &qr.PageRequest)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowNodeTraceLog
func (c *LogController) ShowNodeTraceLog(ctx dotweb.Context) error {
	qr := new(contract.NodeTraceLogQR)
	err := ctx.Bind(qr)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}

	if qr.PageSize <= 0 {
		qr.PageSize = _const.DefaultPageSize
	}

	logService := service.NewLogService()
	result, err := logService.QueryNodeTraceLog(qr.NodeID, &qr.PageRequest)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}
