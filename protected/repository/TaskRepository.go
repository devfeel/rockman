package repository

import (
	"errors"
	"github.com/devfeel/rockman/config"
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
