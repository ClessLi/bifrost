package utils

import (
	"github.com/apsdehal/go-logger"
	"os"
)

var (
	Logger *logger.Logger
	// 日志文件
	Logf    *os.File
	Stdoutf *os.File
)

func InitLogger(stdout *os.File, level logger.LogLevel) {
	var err error
	Logger, err = logger.New("Bifrost", level, stdout)
	if err != nil {
		panic(err)
	}
	Logger.SetFormat("[%{module}] %{time:2006-01-02 15:04:05.000} [%{level}] %{message}\n")
}
