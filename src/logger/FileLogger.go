package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type fileLogger struct {
	filePath string
}

/************************* FileLogger *******************************/

func NewFileLogger(filePath string) *fileLogger {
	if filePath == "" || !fileExists(filePath) {
		filePath = getCurrentDirectory()
	}
	logger := &fileLogger{filePath: filePath}
	return logger
}

func (logger *fileLogger) Print(v ...interface{}) {
	fmt.Println(v...)
}

func (logger *fileLogger) Debug(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "debug")
}
func (logger *fileLogger) Info(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "info")
}
func (logger *fileLogger) Warn(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "warn")
}
func (logger *fileLogger) Error(v ...interface{}) {
	logger.writeLog(fmt.Sprint(v...), "error")
}

func (logger *fileLogger) writeLog(log string, level string) {
	filePath := logger.filePath + "rockman_" + level
	filePath += "_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	logStr := time.Now().Format(defaultFullTimeLayout) + " " + log
	logStr += "\r\n"
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	file, err := os.OpenFile(filePath, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(filePath, err)
		return
	}
	file.WriteString(logStr)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getCurrentDirectory() string {
	return filepath.Clean(filepath.Dir(os.Args[0])) + "/"
}
