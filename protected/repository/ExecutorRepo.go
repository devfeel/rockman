package repository

import (
	"errors"
	"sync"

	"github.com/devfeel/database"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
	"github.com/devfeel/rockman/protected/viewmodel"
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
	sql := "INSERT INTO Task (TaskID,TaskType,IsRun,DueTime,`Interval`,Express,TaskData,TargetType,TargetConfig,DistributeType,Remark)VALUES(?,?,?,?,?,?,?,?,?,?,?);"
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
	sql := "UPDATE Task SET TaskID=?, TaskType = ?, DueTime = ?, `Interval`= ?, Express = ?, TargetType = ?, TargetConfig = ?, Remark = ? WHERE Id = ?;"
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

// IsExistExecutorByTaskId
func (repo *ExecutorRepo) IsExistExecutorByTaskId(taskId string) (bool, error) {
	count, err := repo.Count("SELECT count(1) FROM Task WHERE TaskID=?;", taskId)
	return count > 0, err
}

// QueryExecutors
func (repo *ExecutorRepo) QueryExecutors(qc *viewmodel.ExecutorQC) (*model.PageResult, error) {
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
		err = repo.FindList(&dest, dataSql, params...)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if len(params) != 0 {
		count, err = repo.Count(countSql, params...)
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

// QueryAllExecutors
func (repo *ExecutorRepo) QueryAllExecutors() ([]*model.ExecutorInfo, error) {
	dataSql := "SELECT * FROM Task"
	var dest []*model.ExecutorInfo
	err := repo.FindList(&dest, dataSql)
	if err != nil {
		return nil, err
	}
	return dest, err
}

// WriteExecLog
func (repo *ExecutorRepo) WriteExecLog(log *model.TaskExecLog) (int64, error) {
	sql := "INSERT INTO TaskExecLog(TaskID, NodeID, NodeEndPoint, IsSuccess, StartTime, EndTime, FailureType, FailureCause, CreateTime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return repo.Insert(sql, log.TaskID, log.NodeID, log.NodeEndPoint, log.IsSuccess, log.StartTime, log.EndTime, log.FailureType, log.FailureCause, log.CreateTime)
}

// QueryExecLogs
func (repo *ExecutorRepo) QueryExecLogs(qc *viewmodel.TaskExecLogQC) (*model.PageResult, error) {
	dataSql := "SELECT * FROM TaskExecLog"
	countSql := "SELECT count(1) FROM TaskExecLog"
	if qc.TaskID != "" {
		dataSql += " WHERE TaskID = ?"
		countSql += " WHERE TaskID = ?"
		qc.AddParam(qc.TaskID)
	}
	dataSql += " ORDER BY CreateTime DESC "
	dataSql += qc.GetPageSql()
	params := qc.GetParams()
	var dest []*model.TaskExecLog
	var err error
	if len(params) != 0 {
		err = repo.FindList(&dest, dataSql, params...)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if len(params) != 0 {
		count, err = repo.Count(countSql, params...)
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

// QueryStateLogs
func (repo *ExecutorRepo) QueryStateLogs(qc *viewmodel.TaskStateLogQC) (*model.PageResult, error) {
	dataSql := "SELECT * FROM TaskStateLog"
	countSql := "SELECT count(1) FROM TaskStateLog"
	if qc.TaskID != "" {
		dataSql += " WHERE TaskID = ?"
		countSql += " WHERE TaskID = ?"
		qc.AddParam(qc.TaskID)
	}
	dataSql += " ORDER BY CreateTime DESC "
	dataSql += qc.GetPageSql()
	params := qc.GetParams()
	var dest []*model.TaskExecLog
	var err error
	if len(params) != 0 {
		err = repo.FindList(&dest, dataSql, params...)
	} else {
		err = repo.FindList(&dest, dataSql)
	}
	if err != nil {
		return nil, err
	}

	var count int64
	if len(params) != 0 {
		count, err = repo.Count(countSql, params...)
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
