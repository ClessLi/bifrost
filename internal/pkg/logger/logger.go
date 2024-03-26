package logger

import (
	"github.com/marmotedu/errors"
	"os"
)

type Logger struct {
	initFunc func() error
	flush    func()
}

func (l *Logger) Init() error {
	return l.initFunc()
}

func (l *Logger) Flush() {
	l.flush()
}

func createLogDir(dirpath string) error {
	if info, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if os.IsExist(err) && !info.IsDir() {
		return errors.Errorf("The path '%s' is exist, and it's not a directory", dirpath)
	} else if os.IsExist(err) && info.IsDir() {
		return nil
	} else {
		return err
	}
	return nil
}
