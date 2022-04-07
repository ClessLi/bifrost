package logging

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type loggingWebServerStatusService struct {
	svc svcv1.WebServerStatusService
}

func (l *loggingWebServerStatusService) Get(ctx context.Context) (metrics *v1.Metrics, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Get)
		logF.SetBeginTime(begin)
		defer logF.Result()
		if metrics != nil {
			result, _ := json.Marshal(metrics)
			logF.SetResult(getLimitResult(result))
		}
		logF.SetErr(err)
	}(time.Now().Local())

	return l.svc.Get(ctx)
}

func newWebServerStatusMiddleware(svc svcv1.ServiceFactory) svcv1.WebServerStatusService {
	return &loggingWebServerStatusService{svc: svc.WebServerStatus()}
}
