package repository

import (
	"errors"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
	"sync"
)

const defaultDatabaseID = "demodb"

var defaultRepository *ExecutorRepository
var taskRepositoryLocker *sync.Mutex

func init() {
	taskRepositoryLocker = new(sync.Mutex)
}

type ExecutorRepository struct {
	BaseRepository
}

// GetMessageRepository return MessageRepository which is inited
func GetTaskRepository() *ExecutorRepository {
	//check default repository is init
	if defaultRepository == nil {
		taskRepositoryLocker.Lock()
		defer taskRepositoryLocker.Unlock()
		if defaultRepository == nil {
			defaultRepository = NewTaskRepository()
		}
	}
	return defaultRepository
}

// NewTaskRepository return new MessageRepository
func NewTaskRepository() *ExecutorRepository {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(ExecutorRepository)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

// QueryExecutors
func (repository *ExecutorRepository) QueryExecutors(nodeId string, pageReq *model.PageRequest) (*model.PageResult, error) {
	dataSql := "SELECT * FROM Task"
	countSql := "SELECT count(1) FROM Task"
	if nodeId != "" {
		dataSql += " WHERE TaskID = ?"
		countSql += " WHERE TaskID = ?"
	}
	dataSql += pageReq.GetPageSql()
	var dest []*model.TaskExecLog
	var err error
	if nodeId != "" {
		err = repository.FindList(&dest, dataSql, nodeId)
	} else {
		err = repository.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if nodeId != "" {
		count, err = repository.Count(countSql, nodeId)
	} else {
		count, err = repository.Count(countSql)
	}
	if err != nil {
		return nil, err
	}
	pageResult := new(model.PageResult)
	pageResult.TotalCount = count
	pageResult.PageData = dest
	return pageResult, err
}

// WriteExecLog
func (repository *ExecutorRepository) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repository.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryExecLogs
func (repository *ExecutorRepository) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
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
		err = repository.FindList(&dest, dataSql, taskId)
	} else {
		err = repository.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if taskId != "" {
		count, err = repository.Count(countSql, taskId)
	} else {
		count, err = repository.Count(countSql)
	}
	if err != nil {
		return nil, err
	}
	pageResult := new(model.PageResult)
	pageResult.TotalCount = count
	pageResult.PageData = dest
	return pageResult, err
}
