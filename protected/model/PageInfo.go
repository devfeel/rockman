package model

import (
	"strconv"
)

type (
	PageResult struct {
		TotalCount int64
		PageData   interface{}
	}

	PageRequest struct {
		PageIndex int64
		PageSize  int64
		params    []interface{}
	}
)

func (page *PageRequest) GetSkip() int64 {
	if page.PageIndex <= 0 {
		page.PageIndex = 1
	}
	return (page.PageIndex - 1) * page.PageSize
}

func (page *PageRequest) GetLimit() int64 {
	return page.PageSize
}

func (page *PageRequest) GetPageSql() string {
	return " limit " + strconv.FormatInt(page.GetSkip(), 10) + "," + strconv.FormatInt(page.GetLimit(), 10)
}

func (page *PageRequest) AddParam(param interface{}) *PageRequest {
	page.params = append(page.params, param)
	return page
}

func (page *PageRequest) GetParams() []interface{} {
	return page.params
}
