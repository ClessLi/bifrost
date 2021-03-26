package utils

import (
	"github.com/apsdehal/go-logger"
	"os"
)

func init() {
	var err error
	Logger, err = logger.New("init", logger.DebugLevel, os.Stdout)
	if err != nil {
		panic(err)
	}
	Logger.SetFormat("[%{module}] %{time:2006-01-02 15:04:05.000} [%{level}] %{message}\n")
}
