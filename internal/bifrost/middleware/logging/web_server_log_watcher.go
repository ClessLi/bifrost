package logging

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	"time"
)

type loggingWebServerLogWatcherService struct {
	svc svcv1.WebServerLogWatcherService
}

func (l *loggingWebServerLogWatcherService) Watch(ctx context.Context, request *v1.WebServerLogWatchRequest) (wslog *v1.WebServerLog, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Watch)
		logF.SetBeginTime(begin)
		defer logF.Result()
		if wslog != nil {
			logF.SetResult("Watching web server log...")
		}
		logF.SetErr(err)
	}(time.Now().Local())
	return l.svc.Watch(ctx, request)
}

func newWebServerLogWatcherMiddleware(svc svcv1.ServiceFactory) svcv1.WebServerLogWatcherService {
	return &loggingWebServerLogWatcherService{svc: svc.WebServerLogWatcher()}
}
