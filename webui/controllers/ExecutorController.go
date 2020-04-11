package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman-webui/src/protected/viewModel/task"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
<<<<<<< HEAD
	"github.com/devfeel/rockman/protected/viewmodel"
	"github.com/devfeel/rockman/runtime/executor"
=======
>>>>>>> master
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
<<<<<<< HEAD
	taskService := service.NewExecutorService()
	if model.ID > 0 {
		result = taskService.UpdateExecutor(model)
		if !result.IsSuccess() {
			return ctx.WriteJson(FailedResponse(result.RetCode, "AddExecutor failed: "+result.Message()))
		} else {
			if !model.IsRun {
				//TODO stop executor stop

			}
		}
	} else {
		result = taskService.AddExecutor(model)
		if !result.IsSuccess() {
			return ctx.WriteJson(FailedResponse(result.RetCode, "AddExecutor failed: "+result.Message()))
		}

=======

	result := c.executorService.AddExecutor(model)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(result.RetCode, "AddExecutor failed: "+result.Message()))
>>>>>>> master
	}
	if model.IsRun {
		// submit executor to leader node
		submit := new(core.ExecutorInfo)
		submit.TaskConfig = model.TaskConfig()
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
	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecutorById(model.ID)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "query failed: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

func (c *ExecutorController) DeleteById(ctx dotweb.Context) error {
	model := task.ExecutorInfo{}
	err := ctx.Bind(&model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1001, "parameter bind failed: "+err.Error()))
	}
	result := service.NewExecutorService().RemoveExecutor(model.ID)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(result.RetCode, "RemoveExecutor failed: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(nil))
}

// ShowExecutors
func (c *ExecutorController) ShowExecutors(ctx dotweb.Context) error {
	// nodeId := ctx.QueryString("node")
	// pageIndex := ctx.QueryInt64("pageindex")
	// pageSize := ctx.QueryInt64("pagesize")
	// pageReq := new(model.PageRequest)
	// pageReq.PageIndex = pageIndex
	// pageReq.PageSize = pageSize

	qc := new(viewmodel.ExecutorQC)
	//自动组装参数
	err := ctx.Bind(qc)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}
	if qc.PageSize <= 0 {
		qc.PageSize = _const.DefaultPageSize
	}
	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecutors(qc)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// ShowExecLogs
func (c *ExecutorController) ShowExecLogs(ctx dotweb.Context) error {
	qc := new(viewmodel.TaskExecLogQC)
	//自动组装参数
	err := ctx.Bind(qc)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}

	if qc.PageSize <= 0 {
		qc.PageSize = _const.DefaultPageSize
	}

	taskService := service.NewExecutorService()
	result, err := taskService.QueryExecLogs(qc)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}

// QueryStateLogs
func (c *ExecutorController) QueryStateLogs(ctx dotweb.Context) error {
	qc := new(viewmodel.TaskStateLogQC)
	//自动组装参数
	err := ctx.Bind(qc)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}

	if qc.PageSize <= 0 {
		qc.PageSize = _const.DefaultPageSize
	}

<<<<<<< HEAD
	taskService := service.NewExecutorService()
	result, err := taskService.QueryStateLogs(qc)
=======
	logService := service.NewLogService()
	result, err := logService.QueryExecLogs(taskId, pageReq)
>>>>>>> master
	if err != nil {
		return ctx.WriteJson(FailedResponse(-2001, "Query error: "+err.Error()))
	}
	return ctx.WriteJson(SuccessResponse(result))
}
