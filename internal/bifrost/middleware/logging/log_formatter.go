package logging

import (
	"context"
	"github.com/ClessLi/bifrost/internal/bifrost/middleware/utils"
	"reflect"
	"runtime"
	"time"
)

type logFormatter struct {
	initInfos []interface{}
	result    interface{}
	err       error
	begin     time.Time
}

func (l *logFormatter) AddInfos(infos ...interface{}) {
	l.initInfos = append(l.initInfos, infos...)
}

func (l *logFormatter) Log(infos ...interface{}) {
	logger.Log(append(l.initInfos, infos...)...)
}

func (l *logFormatter) SetResult(result interface{}) {
	l.result = result
}

func (l *logFormatter) SetErr(err error) {
	l.err = err
}

func (l *logFormatter) Result() {
	var infos []interface{}
	if l.result != nil {
		infos = append(infos, "result", l.result)
	}
	if l.err != nil {
		infos = append(infos, "error", l.err)
	}
	l.Log(infos...)
}

func newLogFormatter(ctx context.Context, receiver, method interface{}, begin time.Time) *logFormatter {
	return &logFormatter{
		initInfos: []interface{}{
			//"receiver", reflect.TypeOf(receiver).Kind().String(),
			//"method", reflect.TypeOf(method),
			"method", runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(),
			"clientIp", utils.GetClientIP(ctx),
			"authn", utils.GetAuthnInfo(ctx),
		},
		begin: begin,
	}
}
