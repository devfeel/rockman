package executor

import (
	"errors"
	"sync"

	"github.com/devfeel/database"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/repository"
	"github.com/devfeel/rockman/protected/viewmodel"
)

var defaultRepository *ExecutorRepository
var repositoryLocker *sync.Mutex

func init() {
	repositoryLocker = new(sync.Mutex)
}

type ExecutorRepository struct {
	repository.BaseRepository
}

// GetRepository return ExecutorRepository which is init
func GetRepository() *ExecutorRepository {
	//check default repository is init
	if defaultRepository == nil {
		repositoryLocker.Lock()
		defer repositoryLocker.Unlock()
		if defaultRepository == nil {
			defaultRepository = NewRepository()
		}
	}
	return defaultRepository
}

// NewRepository return new ExecutorRepository
func NewRepository() *ExecutorRepository {
	if config.CurrentProfile.Global.DataBaseConnectString == "" {
		err := errors.New("no config database config")
		panic(err)
	}
	repository := new(ExecutorRepository)
	repository.Init(config.CurrentProfile.Global.DataBaseConnectString)
	repository.InitLogger()
	return repository
}

// InsertOnce
func (repository *ExecutorRepository) InsertOnce(model *model.ExecutorInfo) error {
	sql := "INSERT INTO Task (TaskID,TaskType,IsRun,DueTime,`Interval`,Express,TaskData,TargetType,TargetConfig,DistributeType,Remark)VALUES(?,?,?,?,?,?,?,?,?,?,?);"
	n, err := repository.Insert(sql,
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
func (repository *ExecutorRepository) UpdateOnce(model *model.ExecutorInfo) error {
	sql := "UPDATE Task SET TaskID=?, TaskType = ?, DueTime = ?, `Interval`= ?, Express = ?, TargetType = ?, TargetConfig = ?, Remark = ? WHERE Id = ?;"
	n, err := repository.Update(sql,
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
func (repository *ExecutorRepository) DeleteOnce(id int64) error {
	n, err := repository.Delete("DELETE FROM Task WHERE Id=?;", id)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}
	return nil
}

// GetExecutorById
func (repository *ExecutorRepository) GetExecutorById(id int64) (*model.ExecutorInfo, error) {
	result := &model.ExecutorInfo{}
	err := repository.FindOne(result, "SELECT * FROM Task WHERE Id=?;", id)
	return result, err
}

// GetExecutorByTaskId
func (repository *ExecutorRepository) GetExecutorByTaskId(taskId string) (*model.ExecutorInfo, error) {
	result := &model.ExecutorInfo{}
	err := repository.FindOne(result, "SELECT * FROM Task WHERE TaskID=?;", taskId)
	return result, err
}

// IsExistExecutorByTaskId
func (repository *ExecutorRepository) IsExistExecutorByTaskId(taskId string) (bool, error) {
	count, err := repository.Count("SELECT count(1) FROM Task WHERE TaskID=?;", taskId)
	return count > 0, err
}

// QueryExecutors
func (repository *ExecutorRepository) QueryExecutors(qc *viewmodel.ExecutorQC) (*model.PageResult, error) {
	dataSql := "SELECT * FROM Task"
	countSql := "SELECT count(1) FROM Task"
	if qc.NodeID != "" {
		dataSql += " WHERE NodeID = ?"
		countSql += " WHERE NodeID = ?"
		qc.AddParam(qc.NodeID)
	}
	dataSql += " ORDER BY CreateTime DESC "
	dataSql += qc.GetPageSql()
	params := qc.GetParams()
	var dest []*model.ExecutorInfo
	var err error
	if len(params) != 0 {
		err = repository.FindList(&dest, dataSql, params...)
	} else {
		err = repository.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if len(params) != 0 {
		count, err = repository.Count(countSql, params...)
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
