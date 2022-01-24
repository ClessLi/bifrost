package logging

import (
	"context"
	"encoding/json"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	"time"
)

type loggingWebServerStatisticsService struct {
	svc svcv1.WebServerStatisticsService
}

func (l *loggingWebServerStatisticsService) Get(ctx context.Context, servername *v1.ServerName) (s *v1.Statistics, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Get)
		logF.SetBeginTime(begin)
		defer logF.Result()
		logF.AddInfos(
			"request server name", servername.Name,
		)
		if s != nil {
			result, _ := json.Marshal(s)
			logF.SetResult(getLimitResult(result))
		}

	}(time.Now().Local())
	return l.svc.Get(ctx, servername)
}

func newWebServerStatisticsMiddleware(svc svcv1.ServiceFactory) svcv1.WebServerStatisticsService {
	return &loggingWebServerStatisticsService{svc: svc.WebServerStatistics()}
}
