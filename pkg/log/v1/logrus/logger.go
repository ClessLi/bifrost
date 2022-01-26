package logrus

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// NewLogger create a logrus logger, add hook to it and return it.
func NewLogger(zapLogger *zap.Logger) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logger.AddHook(newHook(zapLogger))

	return logger
}
