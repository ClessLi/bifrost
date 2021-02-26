package logging

import (
	"context"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

type loggingViewer struct {
	viewer service.Viewer
	logger log.Logger
}

func (v loggingViewer) View(requestInfo service.ViewRequestInfo) (responseInfo service.ViewResponseInfo) {
	ip, err := getClientIP(requestInfo.Context())
	defer func(begin time.Time) {
		n := 100
		if responseInfo != nil {
			data := responseInfo.Bytes()
			if data != nil {
				if len(data) < n {
					n = len(data)
				}
			} else {
				data = []byte("")
				n = 0
			}
			v.logger.Log(
				"functions", "View",
				"requestType", requestInfo.GetRequestType(),
				"clientIp", ip,
				"token", requestInfo.GetToken(),
				"webServerName", responseInfo.GetServerName(),
				"result", string(data[:n])+"...",
				"error", responseInfo.Error(),
				"took", time.Since(begin),
			)
			return
		}
		v.logger.Log(
			"functions", "View",
			"requestType", requestInfo.GetRequestType(),
			"clientIp", ip,
			"token", requestInfo.GetToken(),
			"webServerName", requestInfo.GetServerName(),
			"responseInfo", responseInfo,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now().Local())
	if err != nil {
		return service.NewViewResponseInfo(requestInfo.GetServerName(), nil, err)
	}
	responseInfo = v.viewer.View(requestInfo)
	return
}

func loggingViewerMiddleware(logger log.Logger) service.ViewerMiddleware {
	return func(next service.Viewer) service.Viewer {
		return loggingViewer{
			viewer: next,
			logger: logger,
		}
	}
}

type loggingUpdater struct {
	updater service.Updater
	logger  log.Logger
}

func (u loggingUpdater) Update(requestInfo service.UpdateRequestInfo) (responseInfo service.UpdateResponseInfo) {
	ip, err := getClientIP(requestInfo.Context())
	defer func(begin time.Time) {
		if responseInfo != nil {
			u.logger.Log(
				"functions", "Update",
				"requestType", requestInfo.GetRequestType(),
				"clientIp", ip,
				"token", requestInfo.GetToken(),
				"webServerName", responseInfo.GetServerName(),
				"error", responseInfo.Error(),
				"took", time.Since(begin),
			)
			return
		}
		u.logger.Log(
			"functions", "Update",
			"requestType", requestInfo.GetRequestType(),
			"clientIp", ip,
			"token", requestInfo.GetToken(),
			"webServerName", requestInfo.GetServerName(),
			"responseInfo", responseInfo,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now().Local())
	if err != nil {
		return service.NewUpdateResponseInfo(requestInfo.GetServerName(), err)
	}
	responseInfo = u.updater.Update(requestInfo)
	return
}

func loggingUpdaterMiddleware(logger log.Logger) service.UpdaterMiddleware {
	return func(next service.Updater) service.Updater {
		return loggingUpdater{
			updater: next,
			logger:  logger,
		}
	}
}

type loggingWatcher struct {
	watcher service.Watcher
	logger  log.Logger
}

func (w loggingWatcher) Watch(requestInfo service.WatchRequestInfo) (responseInfo service.WatchResponseInfo) {
	ip, err := getClientIP(requestInfo.Context())
	defer func(begin time.Time) {
		if responseInfo != nil {
			w.logger.Log(
				"functions", "Watch",
				"requestType", requestInfo.GetRequestType(),
				"clientIp", ip,
				"token", requestInfo.GetToken(),
				"webServerName", responseInfo.GetServerName(),
				"error", responseInfo.Error(),
				"took", time.Since(begin),
			)
			return
		}
		w.logger.Log(
			"functions", "Watch",
			"requestType", requestInfo.GetRequestType(),
			"clientIp", ip,
			"token", requestInfo.GetToken(),
			"webServerName", requestInfo.GetServerName(),
			"responseInfo", responseInfo,
			"error", err,
			"took", time.Since(begin),
		)

	}(time.Now().Local())
	if err != nil {
		return service.NewWatchResponseInfo(requestInfo.GetServerName(), nil, nil, nil, err)
	}
	responseInfo = w.watcher.Watch(requestInfo)
	return
}

func loggingWatcherMiddleware(logger log.Logger) service.WatcherMiddleware {
	return func(next service.Watcher) service.Watcher {
		return loggingWatcher{
			watcher: next,
			logger:  logger,
		}
	}
}

// loggingMiddleware Make a new type
// that contains BifrostService interface and logger instance
type loggingMiddleware struct {
	viewer  service.Viewer
	updater service.Updater
	watcher service.Watcher
	logger  log.Logger
}

func (lmw loggingMiddleware) Viewer() service.Viewer {
	return lmw.viewer
}

func (lmw loggingMiddleware) Updater() service.Updater {
	return lmw.updater
}

func (lmw loggingMiddleware) Watcher() service.Watcher {
	return lmw.watcher
}

func (lmw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lmw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	result = true
	return
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{
			viewer:  loggingViewerMiddleware(logger)(next.Viewer()),
			updater: loggingUpdaterMiddleware(logger)(next.Updater()),
			watcher: loggingWatcherMiddleware(logger)(next.Watcher()),
			logger:  logger,
		}
	}
}

func getClientIP(ctx context.Context) (ip string, err error) {
	//md, ok := metadata.FromIncomingContext(ctx)
	pr, ok := peer.FromContext(ctx)
	if !ok {
		err = fmt.Errorf("getClientIP, invoke FromContext() failed")
		//err = fmt.Errorf("getClientIP, invoke FromIncomingContext() failed")
		return
	}
	if pr.Addr == net.Addr(nil) {
		err = fmt.Errorf("getClientIP, peer.Addr is nil")
		return
	}
	//fmt.Println(md)
	//ips := md.Get("x-real-ip")
	//if len(ips) == 0 {
	//	err = fmt.Errorf("get real ip failed")
	//	return
	//}
	//ip = ips[0]
	ip = pr.Addr.String()
	return ip, nil
}
