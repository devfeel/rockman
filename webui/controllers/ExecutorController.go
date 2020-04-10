package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/runtime/executor"
	_const "github.com/devfeel/rockman/webui/const"
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
		submit.TaskConfig = getTaskConfig(model)
		if submit.TaskConfig.TargetConfig == nil {
			return ctx.WriteJson(FailedResponse(-1101, "Submit.TaskConfig.TargetConfig is nil"))
		}
		submit.DistributeType = model.DistributeType
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
	nodeId := ctx.QueryString("node")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}
	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecutors(nodeId, pageReq)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecLogs
func (c *ExecutorController) ShowExecLogs(ctx dotweb.Context) error {
	taskId := ctx.QueryString("task")
	pageIndex := ctx.QueryInt64("pageindex")
	pageSize := ctx.QueryInt64("pagesize")
	pageReq := new(model.PageRequest)
	pageReq.PageIndex = pageIndex
	pageReq.PageSize = pageSize

	if pageReq.PageSize <= 0 {
		pageReq.PageSize = _const.DefaultPageSize
	}

	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecLogs(taskId, pageReq)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

func getTaskConfig(model *model.ExecutorInfo) *core.TaskConfig {
	conf := &core.TaskConfig{}
	conf.TaskID = model.TaskID
	conf.TaskType = model.TaskType
	conf.TargetType = model.TargetType
	conf.IsRun = model.IsRun
	conf.DueTime = model.DueTime
	conf.Interval = model.Interval
	conf.Express = model.Express
	conf.TaskData = model.TaskData
	conf.HAFlag = true
	if model.TargetType == executor.TargetType_Http {
		conf.TargetConfig = model.RealTargetConfig.(*executor.HttpConfig)
	}
	if model.TargetType == executor.TargetType_GoSo {
		conf.TargetConfig = model.RealTargetConfig.(*executor.GoConfig)
	}
	if model.TargetType == executor.TargetType_Shell {
		conf.TargetConfig = model.RealTargetConfig.(*executor.ShellConfig)
	}
	return conf
}
