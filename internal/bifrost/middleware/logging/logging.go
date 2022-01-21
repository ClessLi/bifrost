package logging

import (
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	kitlog "github.com/go-kit/kit/log"
	"sync"
)

var logger kitlog.Logger
var limit int
var once = sync.Once{}

type loggingService struct {
	svc svcv1.ServiceFactory
}

func (l *loggingService) WebServerConfig() svcv1.WebServerConfigService {
	return newWebServerConfigMiddleware(l.svc)
}

func (l *loggingService) WebServerStatistics() svcv1.WebServerStatisticsService {
	return newWebServerStatisticsMiddleware(l.svc)
}

func (l *loggingService) WebServerStatus() svcv1.WebServerStatusService {
	return newWebServerStatusMiddleware(l.svc)
}

func (l *loggingService) WebServerLogWatcher() svcv1.WebServerLogWatcherService {
	return newWebServerLogWatcherMiddleware(l.svc)
}

func New(svc svcv1.ServiceFactory) svcv1.ServiceFactory {
	once.Do(func() {
		logger = log.K()
		limit = 100
	})
	return &loggingService{svc: svc}
}

func getLimitResult(result []byte) string {
	if len(result)-1 > limit+3 {
		result = append(result[:limit], '.', '.', '.')
	}
	return string(result)
}
