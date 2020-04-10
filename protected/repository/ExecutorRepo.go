package repository

import (
	"errors"
	"github.com/devfeel/database"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
	"sync"
)

var executorRepo *ExecutorRepo
var executorLocker *sync.Mutex

func init() {
	executorLocker = new(sync.Mutex)
}

type ExecutorRepo struct {
	BaseRepository
}

// GetRepository return ExecutorRepository which is init
func GetExecutorRepo() *ExecutorRepo {
	//check default repository is init
	if executorRepo == nil {
		executorLocker.Lock()
		defer executorLocker.Unlock()
		if executorRepo == nil {
			executorRepo = NewExecutorRepo()
		}
	}
	return executorRepo
}

// NewExecutorRepo return new ExecutorRepo
func NewExecutorRepo() *ExecutorRepo {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(ExecutorRepo)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

// InsertOnce
func (repo *ExecutorRepo) InsertOnce(model *model.ExecutorInfo) error {
	sql := "INSERT INTO Task (TaskID,TaskType,IsRun,DueTime,Interval,Express,TaskData,TargetType,TargetConfig,DistributeType,Remark)VALUES(?,?,?,?,?,?,?,?,?,?,?);"
	n, err := repo.Insert(sql,
		model.TaskID, model.TaskType, 0, model.DueTime, model.Interval,
		model.Express, "", model.TargetType, model.TargetConfig,
		model.DistributeType, model.Remark)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}

	return nil
}

// UpdateOnce
func (repo *ExecutorRepo) UpdateOnce(model *model.ExecutorInfo) error {
	sql := "UPDATE Task SET TaskID=?, TaskType = ?, DueTime = ?, Interval= ?, Express = ?, TargetType = ?, TargetConfig = ?, Remark = ? WHERE Id = ?;"
	n, err := repo.Update(sql,
		model.TaskID,
		model.TaskType, model.DueTime, model.Interval, model.Express,
		model.TargetType, model.TargetConfig, model.Remark, model.ID)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}
	return nil
}

// DeleteOnce
func (repo *ExecutorRepo) DeleteOnce(id int64) error {
	n, err := repo.Delete("DELETE FROM Task WHERE Id=?;", id)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}
	return nil
}

// GetExecutorById
func (repo *ExecutorRepo) GetExecutorById(id int64) (*model.ExecutorInfo, error) {
	result := &model.ExecutorInfo{}
	err := repo.FindOne(result, "SELECT * FROM Task WHERE Id=?;", id)
	return result, err
}

// GetExecutorByTaskId
func (repo *ExecutorRepo) GetExecutorByTaskId(taskId string) (*model.ExecutorInfo, error) {
	result := &model.ExecutorInfo{}
	err := repo.FindOne(result, "SELECT * FROM Task WHERE TaskID=?;", taskId)
	return result, err
}

// QueryExecutors
func (repo *ExecutorRepo) QueryExecutors(nodeId string, pageReq *model.PageRequest) (*model.PageResult, error) {
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
		err = repo.FindList(&dest, dataSql, nodeId)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if nodeId != "" {
		count, err = repo.Count(countSql, nodeId)
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

// WriteExecLog
func (repo *ExecutorRepo) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryExecLogs
func (repo *ExecutorRepo) QueryExecLogs(taskId string, pageReq *model.PageRequest) (*model.PageResult, error) {
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
func (repo *ExecutorRepo) WriteNodeTraceLog(log *model.NodeTraceLog) (int64, error) {
	sql := "INSERT INTO NodeTraceLog(NodeID, NodeEndPoint, IsLeader, IsMaster, IsWorker, Event, IsSuccess, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.NodeID, log.NodeEndPoint, log.IsLeader, log.IsMaster, log.IsWorker, log.Event, log.IsSuccess, log.FailureType, log.FailureCause, log.CreateTime)
}
