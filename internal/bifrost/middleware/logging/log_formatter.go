package logging

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/yongPhone/bifrost/internal/bifrost/middleware/utils"
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
	if !reflect.DeepEqual(l.begin, time.Time{}) {
		infos = append(infos, "took", time.Since(l.begin))
	}
	logger.Log(append(l.initInfos, infos...)...) //nolint:errcheck
}

func (l *logFormatter) SetResult(result interface{}) {
	l.result = result
}

func (l *logFormatter) SetErr(err error) {
	l.err = err
}

func (l *logFormatter) SetBeginTime(begin time.Time) {
	l.begin = begin
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

func newLogFormatter(ctx context.Context, method interface{}) *logFormatter {
	return &logFormatter{
		initInfos: []interface{}{
			"method", runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(),
			"clientIp", utils.GetClientIP(ctx),
			"authn", utils.GetAuthnInfo(ctx),
		},
		begin: time.Time{},
	}
}
