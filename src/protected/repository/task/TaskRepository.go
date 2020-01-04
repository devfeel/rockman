package task

import (
	"errors"
	"github.com/devfeel/rockman/src/config"
	"github.com/devfeel/rockman/src/protected/repository"
	"sync"
)

const defaultDatabaseID = "demodb"

var defaultTaskRepository *TaskRepository
var taskRepositoryLocker *sync.Mutex

func init() {
	taskRepositoryLocker = new(sync.Mutex)
}

type TaskRepository struct {
	repository.BaseRepository
}

// GetMessageRepository return MessageRepository which is inited
func GetTaskRepository() *TaskRepository {
	//check default repository is init
	if defaultTaskRepository == nil {
		taskRepositoryLocker.Lock()
		defer taskRepositoryLocker.Unlock()
		if defaultTaskRepository == nil {
			defaultTaskRepository = NewDemoRepository()
		}
	}
	return defaultTaskRepository
}

// NewMessageRepository return new MessageRepository
func NewDemoRepository() *TaskRepository {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(TaskRepository)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

func (repository *TaskRepository) QueryTasksByNodeID(dest interface{}, taskId string) error {
	sql := "SELECT * FROM [Tasks] WHERE NodeID = ? "
	return repository.FindList(dest, sql, taskId)
}
