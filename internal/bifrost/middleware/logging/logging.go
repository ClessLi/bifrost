package logging

import (
	"sync"

	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	kitlog "github.com/go-kit/kit/log"
)

var (
	logger kitlog.Logger
	limit  int
	once   = sync.Once{}
)

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

func (l *loggingService) WebServerBinCMD() svcv1.WebServerBinCMDService {
	return newWebServerBinCMDMiddleware(l.svc)
}

func New(svc svcv1.ServiceFactory) svcv1.ServiceFactory {
	once.Do(func() {
		logger = logV1.K()
		limit = 100
	})

	return &loggingService{svc: svc}
}

func getLimitResult(result []byte) string {
	var formattedRet []byte
	if len(result)-1 > limit+3 {
		formattedRet = result[:limit]
	} else {
		formattedRet = result
	}

	return string(formattedRet) + "..."
}
