package core

type JsonResult struct {
	RetCode int
	RetMsg  string
	Message interface{}
}

func (r *JsonResult) CorrectCode() int {
	return 0
}
