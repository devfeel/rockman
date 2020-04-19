package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	_const "github.com/devfeel/rockman/webui/const"
	"github.com/devfeel/rockman/webui/contract"
)

type ExecutorController struct {
	executorService *service.ExecutorService
}

func NewExecutorController() *ExecutorController {
	return &ExecutorController{
		executorService: service.NewExecutorService(),
	}
}

// SaveExecutor
func (c *ExecutorController) SaveExecutor(ctx dotweb.Context) error {
	model := &model.ExecutorInfo{}
	err := ctx.Bind(model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}

	result := c.executorService.AddExecutor(model)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(result.RetCode, "AddExecutor failed: "+result.Message()))
	}
	if model.IsRun {
		// submit executor to leader node
		submit := new(core.ExecutorInfo)
		submit.TaskConfig = model.TaskConfig()
		if submit.TaskConfig.TargetConfig == nil {
			return ctx.WriteJson(FailedResponse(-1101, "Submit.TaskConfig.TargetConfig is nil"))
		}
		leader := getLeader(ctx)
		if leader == "" {
			return ctx.WriteJson(FailedResponse(-1102, "Leader is nil"))
		}
		// submit to rpc
		err, reply := GetRpcClient(leader).CallSubmitExecutor(submit)
		if err != nil {
			return ctx.WriteJson(FailedResponse(-1201, "CallSubmitExecutor error: "+err.Error()))
		} else {
			if reply.IsSuccess() {
				//TODO update db IsSubmit = true
			} else {
				return ctx.WriteJson(FailedResponse(-1202, "CallSubmitExecutor failed: "+reply.RetMsg))
			}
		}
	}
	return ctx.WriteJson(SuccessResponse(nil))
}

// UpdateExecutor
func (c *ExecutorController) UpdateExecutor(ctx dotweb.Context) error {
	model := &model.ExecutorInfo{}
	err := ctx.Bind(model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}

	dbExecInfo, err := c.executorService.QueryExecutorById(model.ID)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1003, "query task error:"+err.Error()))
	}
	if dbExecInfo == nil {
		return ctx.WriteJson(FailedResponse(-1004, "not exists this task"))
	}

	result := c.executorService.UpdateExecutor(model)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(result.RetCode, "UpdateExecutor failed: "+result.Message()))
	} else {
		leader := getLeader(ctx)
		if leader == "" {
			return ctx.WriteJson(FailedResponse(-1101, "Leader is nil"))
		}
		if dbExecInfo.IsRun && !model.IsRun {
			err, reply := GetRpcClient(leader).CallSubmitStopExecutor(model.TaskID)
			if err != nil {
				return ctx.WriteJson(FailedResponse(-1201, "CallSubmitStopExecutor error: "+err.Error()))
			} else {
				if reply.IsSuccess() {
					//TODO log something
				} else {
					return ctx.WriteJson(FailedResponse(-1202, "CallSubmitStopExecutor failed: "+reply.RetMsg))
				}
			}
		}
		if !dbExecInfo.IsRun && model.IsRun {
			err, reply := GetRpcClient(leader).CallSubmitStartExecutor(model.TaskID)
			if err != nil {
				return ctx.WriteJson(FailedResponse(-1201, "CallSubmitStartExecutor error: "+err.Error()))
			} else {
				if reply.IsSuccess() {
					//TODO log something
				} else {
					return ctx.WriteJson(FailedResponse(-1202, "CallSubmitStartExecutor failed: "+reply.RetMsg))
				}
			}
		}
	}
	return ctx.WriteJson(SuccessResponse(nil))
}

// QueryById
func (c *ExecutorController) QueryById(ctx dotweb.Context) error {
	model := model.ExecutorInfo{}
	err := ctx.Bind(&model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}
	result, err := c.executorService.QueryExecutorById(model.ID)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "query failed: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecutors
func (c *ExecutorController) ShowExecutors(ctx dotweb.Context) error {
	qr := new(contract.ExecutorQR)
	err := ctx.Bind(qr)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}
	if qr.PageSize <= 0 {
		qr.PageSize = _const.DefaultPageSize
	}

	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecutors(&qr.PageRequest)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecLogs
func (c *ExecutorController) ShowExecLogs(ctx dotweb.Context) error {
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
func (c *ExecutorController) ShowStateLog(ctx dotweb.Context) error {
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
