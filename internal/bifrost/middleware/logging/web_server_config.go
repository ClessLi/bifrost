package logging

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type loggingWebServerConfigService struct {
	svc svcv1.WebServerConfigService
}

func (l loggingWebServerConfigService) GetServerNames(ctx context.Context) (servernames *v1.ServerNames, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.GetServerNames)
		logF.SetBeginTime(begin)
		defer logF.Result()
		if servernames != nil {
			var result []string
			for _, serverName := range *servernames {
				result = append(result, serverName.Name)
			}
			logF.SetResult("ServerNames: " + strings.Join(result, ", "))
		}
		logF.SetErr(err)
	}(time.Now().Local())

	return l.svc.GetServerNames(ctx)
}

func (l loggingWebServerConfigService) Get(
	ctx context.Context,
	servername *v1.ServerName,
) (wsc *v1.WebServerConfig, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Get)
		logF.SetBeginTime(begin)
		defer logF.Result()
		logF.AddInfos(
			"request server name", servername.Name,
		)
		if wsc != nil {
			logF.SetResult(getLimitResult(wsc.JsonData))
		}
		logF.SetErr(err)
	}(time.Now().Local())

	return l.svc.Get(ctx, servername)
}

func (l loggingWebServerConfigService) ConnectivityCheckOfProxiedServers(ctx context.Context, pos *v1.WebServerConfigContextPos) (ctxData *v1.ContextData, err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.ConnectivityCheckOfProxiedServers)
		logF.SetBeginTime(begin)
		defer logF.Result()
		logF.AddInfos(
			"request server name", pos.ServerName.Name,
		)
		if err == nil {
			logF.SetResult(fmt.Sprintf("check connectivity of proxied servers successfully, position: %v, proxy pass data: %s", pos.ContextPos, ctxData.JsonData))
		}
	}(time.Now().Local())

	return l.svc.ConnectivityCheckOfProxiedServers(ctx, pos)
}

func (l loggingWebServerConfigService) Update(ctx context.Context, config *v1.WebServerConfig) (err error) {
	defer func(begin time.Time) {
		logF := newLogFormatter(ctx, l.svc.Update)
		logF.SetBeginTime(begin)
		defer logF.Result()
		logF.AddInfos(
			"request server name", config.ServerName.Name,
		)
		if err == nil {
			logF.SetResult("update web server config succeeded")
		}
		logF.SetErr(err)
	}(time.Now().Local())

	return l.svc.Update(ctx, config)
}

func newWebServerConfigMiddleware(svc svcv1.ServiceFactory) svcv1.WebServerConfigService {
	return &loggingWebServerConfigService{svc: svc.WebServerConfig()}
}
