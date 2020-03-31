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

func (repository *TaskRepository) QueryTasks(dest interface{}) error {
	sql := "SELECT * FROM Task"
	return repository.FindList(dest, sql)
}

func (repository *TaskRepository) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repository.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

func (repository *TaskRepository) QueryLogs(dest interface{}) error {
	sql := "SELECT * FROM TaskExecLog LIMIT 100"
	return repository.FindList(dest, sql)
}
