package repository

import (
	"errors"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
	"sync"
)

const defaultDatabaseID = "demodb"

var defaultTaskRepository *TaskRepository
var taskRepositoryLocker *sync.Mutex

func init() {
	taskRepositoryLocker = new(sync.Mutex)
}

type TaskRepository struct {
	BaseRepository
}

// GetMessageRepository return MessageRepository which is inited
func GetTaskRepository() *TaskRepository {
	//check default repository is init
	if defaultTaskRepository == nil {
		taskRepositoryLocker.Lock()
		defer taskRepositoryLocker.Unlock()
		if defaultTaskRepository == nil {
			defaultTaskRepository = NewTaskRepository()
		}
	}
	return defaultTaskRepository
}

// NewTaskRepository return new MessageRepository
func NewTaskRepository() *TaskRepository {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(TaskRepository)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

func (repository *TaskRepository) QueryTasks() ([]*model.TaskInfo, error) {
	sql := "SELECT * FROM Task"
	var dest []*model.TaskInfo
	var err error
	err = repository.FindList(&dest, sql)
	return dest, err
}

// WriteExecLog
func (repository *TaskRepository) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repository.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryExecLogs
func (repository *TaskRepository) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
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
