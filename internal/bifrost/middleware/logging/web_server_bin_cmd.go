package logging

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type loggingWebServerBinCMDService struct {
	svc svcv1.WebServerBinCMDService
}

func (l *loggingWebServerBinCMDService) Exec(ctx context.Context, request *v1.ExecuteRequest) (response *v1.ExecuteResponse, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Exec)
		logF.SetBeginTime(begin)
		defer logF.Result()
		if response != nil {
			result, _ := json.Marshal(response)
			logF.SetResult(getLimitResult(result))
		}
		logF.SetErr(err)
	}(time.Now().Local())

	return l.svc.Exec(ctx, request)
}

func newWebServerBinCMDMiddleware(svc svcv1.ServiceFactory) svcv1.WebServerBinCMDService {
	return &loggingWebServerBinCMDService{svc: svc.WebServerBinCMD()}
}
