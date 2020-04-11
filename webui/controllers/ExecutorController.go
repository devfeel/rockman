package controllers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman-webui/src/protected/viewModel/task"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/node"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	"github.com/devfeel/rockman/protected/viewmodel"
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
	result, node := getNode(ctx)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(-1001, result.Message()))
	}

	model := &model.ExecutorInfo{}
	err := ctx.Bind(model)
	if err != nil {
		return ctx.WriteJson(FailedResponse(-1002, "parameter bind failed: "+err.Error()))
	}
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

	}

	if model.IsRun {
		// submit executor to leader node
		submit := new(core.ExecutorInfo)
		submit.TaskConfig = getTaskConfig(model)
		if submit.TaskConfig.TargetConfig == nil {
			return ctx.WriteJson(FailedResponse(-1101, "Submit.TaskConfig.TargetConfig is nil"))
		}
		submit.DistributeType = model.DistributeType
		// submit to rpc
		err, reply := node.Cluster.GetLeaderRpcClient().CallSubmitExecutor(submit)
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
	result, node := getNode(ctx)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(-1001, result.Message()))
	}

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

	result = c.executorService.UpdateExecutor(model)
	if !result.IsSuccess() {
		return ctx.WriteJson(FailedResponse(result.RetCode, "UpdateExecutor failed: "+result.Message()))
	} else {
		if dbExecInfo.IsRun && !model.IsRun {
			err, reply := node.Cluster.GetLeaderRpcClient().CallSubmitStopExecutor(model.TaskID)
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
			err, reply := node.Cluster.GetLeaderRpcClient().CallSubmitStartExecutor(model.TaskID)
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

	taskService := service.NewExecutorService()
	result, err := taskService.QueryStateLogs(qc)
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

func getNode(ctx dotweb.Context) (*core.Result, *node.Node) {
	nodeItem, exists := ctx.AppItems().Get(_const.ItemKey_Node)
	if !exists {
		return core.FailedResult(-1001, "not exists node in app items"), nil
	}
	node, ok := nodeItem.(*node.Node)
	if !ok {
		return core.FailedResult(-1002, "not exists correct node in app items"), nil
	}
	return core.SuccessResult(), node
}
