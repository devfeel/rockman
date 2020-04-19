package repository

import (
	"errors"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
)

type LogRepo struct {
	BaseRepository
}

// NewLogRepo return new ExecutorRepo
func NewLogRepo() *LogRepo {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(LogRepo)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

// WriteExecLog
func (repo *LogRepo) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryExecLogs
func (repo *LogRepo) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	dataSql := "SELECT * FROM TaskExecLog"
	countSql := "SELECT count(1) FROM TaskExecLog"
	if taskId != "" {
		dataSql += " WHERE TaskID = ?"
		countSql += " WHERE TaskID = ?"
	}
	dataSql += pageReq.GetPageSql()
	var dest []*model.TaskExecLog
	var err error
	if taskId != "" {
		err = repo.FindList(&dest, dataSql, taskId)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if taskId != "" {
		count, err = repo.Count(countSql, taskId)
	} else {
		count, err = repo.Count(countSql)
	}
	if err != nil {
		return nil, err
	}
	pageResult := new(model.PageResult)
	pageResult.TotalCount = count
	pageResult.PageData = dest
	return pageResult, err
}

// WriteNodeTraceLog
func (repo *LogRepo) WriteNodeTraceLog(log *model.NodeTraceLog) (int64, error) {
	sql := "INSERT INTO NodeTraceLog(NodeID, NodeEndPoint, IsLeader, IsMaster, IsWorker, Event, IsSuccess, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.NodeID, log.NodeEndPoint, log.IsLeader, log.IsMaster, log.IsWorker, log.Event, log.IsSuccess, log.FailureType, log.FailureCause, log.CreateTime)
}

// WriteSubmitLog
func (repo *LogRepo) WriteSubmitLog(log *model.TaskSubmitLog) (int64, error) {
	sql := "INSERT INTO TaskSubmitLog(TaskID, NodeID, NodeEndPoint, IsSuccess,  FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryStateLog
func (repo *LogRepo) QueryStateLog(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	dataSql := "SELECT * FROM TaskStateLog"
	countSql := "SELECT count(1) FROM TaskStateLog"
	if taskId != "" {
		dataSql += " WHERE TaskID = ?"
		countSql += " WHERE TaskID = ?"
	}
	dataSql += pageReq.GetPageSql()
	var dest []*model.TaskStateLog
	var err error
	if taskId != "" {
		err = repo.FindList(&dest, dataSql, taskId)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if taskId != "" {
		count, err = repo.Count(countSql, taskId)
	} else {
		count, err = repo.Count(countSql)
	}
	if err != nil {
		return nil, err
	}
	pageResult := new(model.PageResult)
	pageResult.TotalCount = count
	pageResult.PageData = dest
	return pageResult, err
}
