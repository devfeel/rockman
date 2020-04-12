package node

import (
	"errors"
	"fmt"
	"github.com/devfeel/rockman/core"
	"github.com/devfeel/rockman/logger"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/service"
	"github.com/hashicorp/consul/api"
	"strconv"
	"time"
)

func (n *Node) SubmitExecutor(execInfo *core.ExecutorInfo) *core.Result {
	result := n.submitExecutor(execInfo)
	if result.IsSuccess() {
		n.setExecutorChangeFlag()
	}
	return result
}

func (n *Node) SubmitStopExecutor(taskId string) *core.Result {
	result := n.submitStopExecutor(taskId)
	if result.IsSuccess() {
		n.setExecutorChangeFlag()
	}
	return result
}

func (n *Node) SubmitStartExecutor(taskId string) *core.Result {
	result := n.submitStartExecutor(taskId)
	if result.IsSuccess() {
		n.setExecutorChangeFlag()
	}
	return result
}

func (n *Node) submitExecutor(execInfo *core.ExecutorInfo) *core.Result {
	logTitle := "Node SubmitExecutor [" + execInfo.TaskConfig.TaskID + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	if execInfo.Worker != nil {
		if execInfo.Worker.Cluster != n.Cluster.ClusterId {
			logger.Node().Debug(logTitle + "failed, not match cluster [" + execInfo.Worker.Cluster + ", " + n.Cluster.ClusterId + "]")
			return core.FailedResult(-1002, "not match cluster ["+execInfo.Worker.Cluster+", "+n.Cluster.ClusterId+"]")
		}

		endPoint := execInfo.Worker.EndPoint()
		node, exists := n.Cluster.FindNode(endPoint)
		if !exists {
			logger.Node().Debug(logTitle + "failed, can not find node[" + endPoint + "] in cluster")
			return core.FailedResult(-1003, "can not find node["+endPoint+"] in cluster")
		}

		if node.NodeID != execInfo.Worker.NodeID {
			logger.Node().Debug(logTitle + "failed, not match node id [" + execInfo.Worker.NodeID + ", " + node.NodeID + "]")
			return core.FailedResult(-1004, "not match node id ["+execInfo.Worker.NodeID+", "+node.NodeID+"]")
		}
	}

	var err error
	// get low balance worker
	if execInfo.Worker == nil {
		execInfo.Worker, err = n.Cluster.GetLowBalanceWorker()
		if err != nil {
			logger.Node().Error(err, logTitle+"GetLowBalanceWorker error")
			//TODO log submit result to db log
			return core.ErrorResult(err)
		}
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(execInfo.Worker.EndPoint())
	err, reply := rpcClient.CallRegisterExecutor(execInfo.TaskConfig)
	submitLog := &model.TaskSubmitLog{
		TaskID:       execInfo.TaskConfig.TaskID,
		NodeID:       execInfo.Worker.NodeID,
		NodeEndPoint: execInfo.Worker.EndPoint(),
		IsSuccess:    err != nil && reply.IsSuccess(),
	}

	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+execInfo.Worker.EndPoint()+"] error:", err.Error())
		submitLog.FailureType = "error"
		submitLog.FailureCause = err.Error()
		n.logLogic.WriteSubmitLog(submitLog)
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			submitLog.FailureType = "failure"
			submitLog.FailureCause = reply.FailureMessage()
			logger.Node().DebugS(logTitle+"to ["+execInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			n.Cluster.AddExecutor(execInfo)
			submitLog.IsSuccess = true
			logger.Node().Debug(logTitle + "to [" + execInfo.Worker.EndPoint() + "] success.")
		}
		n.logLogic.WriteSubmitLog(submitLog)
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

func (n *Node) submitStopExecutor(taskId string) *core.Result {
	logTitle := "Node SubmitStopExecutor [" + taskId + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	runExecInfo, exists := n.Cluster.FindExecutor(taskId)
	if !exists {
		logger.Node().Debug(logTitle + "failed, can not find executor is running cluster.")
		return core.FailedResult(-1001, "can not find executor is running cluster")
	}

	if !runExecInfo.TaskConfig.IsRun {
		return core.SuccessResult()
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(runExecInfo.Worker.EndPoint())
	err, reply := rpcClient.CallStopExecutor(taskId)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] error:", err.Error())
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			runExecInfo.TaskConfig.IsRun = false
			logger.Node().Debug(logTitle + "to [" + runExecInfo.Worker.EndPoint() + "] success.")
		}
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

func (n *Node) submitStartExecutor(taskId string) *core.Result {
	logTitle := "Node SubmitStartExecutor [" + taskId + "] "
	if !n.IsLeader() {
		logger.Node().Debug(logTitle + "failed, current node is not leader.")
		return core.FailedResult(-1001, "current node is not leader")
	}

	runExecInfo, exists := n.Cluster.FindExecutor(taskId)
	if !exists {
		logger.Node().Debug(logTitle + "failed, can not find executor is running cluster.")
		return core.FailedResult(-1001, "can not find executor is running cluster")
	}

	if runExecInfo.TaskConfig.IsRun {
		return core.SuccessResult()
	}

	//submit executor to the specified worker node
	rpcClient := n.Cluster.GetRpcClient(runExecInfo.Worker.EndPoint())
	err, reply := rpcClient.CallStartExecutor(runExecInfo.TaskConfig.TaskID)
	//TODO log submit result to db log
	if err != nil {
		logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] error:", err.Error())
		return core.ErrorResult(err)
	} else {
		if !reply.IsSuccess() {
			logger.Node().DebugS(logTitle+"to ["+runExecInfo.Worker.EndPoint()+"] failed, result:", reply.RetCode)
		} else {
			runExecInfo.TaskConfig.IsRun = true
			logger.Node().Debug(logTitle + "to [" + runExecInfo.Worker.EndPoint() + "] success.")
		}
		return core.NewResult(reply.RetCode, reply.RetMsg, nil)
	}
}

// electionLeader
func (n *Node) electionLeader() {
	logTitle := "Node election leader "
	logger.Node().Debug(logTitle + "begin.")

	doQuery := func() error {
		err := n.Cluster.ElectionLeader(n.NodeInfo().EndPoint())
		if err != nil {
			logger.Node().DebugS(logTitle + "error: " + err.Error() + ", will retry 10 seconds after")
			logger.Node().Error(err, logTitle+"error")
			time.Sleep(time.Second * 10)
			return err
		} else {
			logger.Node().Debug(logTitle + "success with key {" + n.Cluster.LeaderKey + "}")
			n.becomeLeaderRole()
			return nil
		}
	}

	go func() {
		var retryCount int
		limit := n.config.Global.RetryLimit
		for {
			if n.isSTW {
				return
			}
			err := doQuery()
			if err != nil {
				retryCount += 1
				if retryCount > limit {
					err := errors.New("retry count bigger than " + strconv.Itoa(limit) + ", now will stop node.")
					logger.Node().DebugS(logTitle + "error:" + err.Error())
					n.Shutdown()
					return
				}
			} else {
				retryCount = 0
			}
		}
	}()
}

func (n *Node) becomeLeaderRole() {
	logTitle := "Node.becomeLeaderRole "
	logger.Node().Debug(logTitle + "become to leader role")
	n.isLeader = true

	//TODO sync all executors
	n.Cluster.OnNodeOffline = n.onWorkerNodeOffline

}

func (n *Node) removeLeaderRole() {
	logTitle := "Node "
	//TODO do something when become to not leader
	logger.Node().Debug(logTitle + "remove leader role")
	n.Cluster.OnNodeOffline = nil
	n.isLeader = false
}

// initExecutorsFromDB init executors from db
// must check init flag on registry
func (n *Node) prepareExecutorsFromDB() {
	logTitle := "Node initExecutorsFromDB "
	if !n.IsLeader() {
		return
	}
	var successCount, failureCount int
	doQuery := func() {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, logTitle+"throw unhandled error:"+errInfo.Error())
			}
		}()
		execInfos, err := service.NewExecutorService().QueryAllExecutors()
		if err != nil {
			logger.Node().Debug(logTitle + "NewExecutorService error:" + err.Error())
			return
		}
		if execInfos == nil {
			return
		}
		for _, exec := range execInfos {
			submit := new(core.ExecutorInfo)
			submit.TaskConfig = exec.TaskConfig()
			if submit.TaskConfig == nil || submit.TaskConfig.TargetConfig == nil {
				logger.Node().Debug(logTitle + "init submit error: TaskConfig is nil or target config is nil")
				failureCount += 1
				continue
			}
			result := n.SubmitExecutor(submit)
			if result.Error != nil {
				logger.Node().DebugS(logTitle+"SubmitExecutor error:", result.Error.Error())
				failureCount += 1
				//TODO log to db
				continue
			}

			if !result.IsSuccess() {
				logger.Node().DebugS(logTitle + "SubmitExecutor failed, " + result.Message())
				failureCount += 1
				//TODO log to db
			} else {
				logger.Node().DebugS(logTitle + "SubmitExecutor success")
				successCount += 1
				//TODO log to db
			}
		}
	}

	logger.Node().Debug(logTitle + "init begin.")
	flag, err := n.getInitFlag()
	if err != nil {
		logger.Node().Warn(logTitle + "get init flag error:" + err.Error())
	} else {
		if !flag {
			doQuery()
			err := n.setInitFlag()
			if err != nil {
				logger.Node().Warn(logTitle + "set init flag error:" + err.Error())
			}
			logger.Node().Debug(logTitle + "init finish. Success[" + strconv.Itoa(successCount) + "] Failure[" + strconv.Itoa(failureCount) + "]")
		}
	}
}

func (n *Node) watchExecutorChangeFromLeader() {
	lt := "Node.watchExecutorChangeFromLeader "

	doQuery := func() error {
		defer func() {
			if err := recover(); err != nil {
				errInfo := errors.New(fmt.Sprintln(err))
				logger.Cluster().Error(errInfo, lt+"throw unhandled error:"+errInfo.Error())
			}
		}()

		opt := &api.QueryOptions{
			WaitIndex: n.executorFlagLastIndex,
			WaitTime:  time.Minute * 10,
		}
		_, meta, err := n.Registry.Get(getExecutorChangeFlagKey(n.ClusterId()), opt)
		if err != nil {
			return err
		}
		if meta.LastIndex != n.executorFlagLastIndex {
			logger.Cluster().Debug("Cluster.watchLeaderChange: leader changed.")
			leader, err := n.Cluster.GetLeaderInfo()
			if err != nil {
				logger.Node().Error(err, lt+"get leader info error")
				return err
			}
			rpcClient := n.Cluster.GetRpcClient(leader)
			err, reply := rpcClient.CallQueryExecutorInfos("")
			if err != nil {
				logger.Node().Debug(lt + "CallQueryExecutorInfos Error: " + err.Error())
				return err
			}
			if !reply.IsSuccess() {
				logger.Node().Debug(lt + "CallQueryExecutorInfos Failed: " + reply.FailureMessage())
				return errors.New("CallQueryExecutorInfos Failed: " + reply.FailureMessage())
			}
			n.Cluster.ExecutorInfos = reply.Message.(map[string]*core.ExecutorInfo)
			logger.Node().Debug(lt + "CallQueryExecutorInfos Success.")
		}
		return nil
	}
	for {
		doQuery()
	}
}
