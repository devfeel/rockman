package repository

import (
	"errors"
	"github.com/devfeel/database"
	"github.com/devfeel/rockman/config"
	"github.com/devfeel/rockman/protected/model"
)

type ExecutorRepo struct {
	BaseRepository
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

// UpdateSubmitFlag
func (repo *ExecutorRepo) UpdateSubmitFlag(id int64, flag bool) error {
	n, err := repo.Delete("UPDATE Task SET IsSubmitToCluster = ? WHERE Id=?;", id, flag)
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
func (repo *ExecutorRepo) QueryExecutors(pageReq *model.PageRequest) (*model.PageResult, error) {
	dataSql := "SELECT * FROM Task"
	countSql := "SELECT count(1) FROM Task"
	dataSql += " ORDER BY CreateTime DESC " + pageReq.GetPageSql()
	var dest []*model.ExecutorInfo
	err := repo.FindList(&dest, dataSql)
	if err != nil {
		return nil, err
	}
	count, err := repo.Count(countSql)
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

// QueryRunInfo
func (repo *ExecutorRepo) QueryRunInfo(taskId string) (*model.ExecutorRunInfo, error) {
	result := &model.ExecutorRunInfo{}
	err := repo.FindOne(result, "SELECT * FROM ExecutorRunInfo WHERE TaskID=?;", taskId)
	return result, err
}

// InsertRunInfo
func (repo *ExecutorRepo) InsertRunInfo(model *model.ExecutorRunInfo) error {
	sql := "INSERT INTO ExecutorRunInfo (TaskID, NodeID, NodeEndPoint, LastUpdateTime, CreateTime)VALUES(?,?,?,?,?);"
	n, err := repo.Insert(sql,
		model.TaskID, model.NodeID, model.NodeEndPoint, model.LastUpdateTime, model.CreateTime)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}

	return nil
}

// UpdateRunInfo
func (repo *ExecutorRepo) UpdateRunInfo(model *model.ExecutorRunInfo) error {
	sql := "UPDATE ExecutorRunInfo SET NodeID = ?, NodeEndPoint = ?, LastUpdateTime= ? WHERE TaskID = ?;"
	n, err := repo.Update(sql,
		model.NodeID, model.NodeEndPoint, model.LastUpdateTime, model.TaskID)
	if err != nil {
		return err
	}

	if n <= 0 {
		return database.ErrorNoRowsAffected
	}
	return nil
}
