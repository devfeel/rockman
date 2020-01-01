package logger

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
)

type Logger interface {
	Error(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
	Print(v ...interface{})
}
