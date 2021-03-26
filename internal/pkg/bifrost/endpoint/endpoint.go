package endpoint

import (
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

var (
	ErrInvalidRequest = errors.New("request has only one class: Request")
	ErrInvalidService = errors.New("service is invalid or nil")
)

// BifrostEndpoints bifrost gRPC服务Endpoint层结构体
type BifrostEndpoints struct {
	ViewerEndpoint      endpoint.Endpoint
	UpdaterEndpoint     endpoint.Endpoint
	WatcherEndpoint     endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func NewBifrostEndpoints(svc service.Service) BifrostEndpoints {
	return BifrostEndpoints{
		ViewerEndpoint:      MakeViewerEndpoint(svc.Viewer()),
		UpdaterEndpoint:     MakeUpdaterEndpoint(svc.Updater()),
		WatcherEndpoint:     MakeWatcherEndpoint(svc.Watcher()),
		HealthCheckEndpoint: MakeHealthCheckEndpoint(svc),
	}
}

func MakeViewerEndpoint(viewer service.Viewer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//utils.Logger.Debug("request to bifrost viewer endpoint")
		if req, ok := request.(*ViewRequestInfo); ok {
			svcRequestInfo := service.NewViewRequestInfo(ctx, req.ViewType, req.ServerName, req.Token)
			svcResponseInfo := viewer.View(svcRequestInfo)
			viewResponseInfo := newViewResponseInfo(svcResponseInfo)
			return viewResponseInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

func MakeUpdaterEndpoint(updater service.Updater) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//utils.Logger.Debug("request to bifrost updater endpoint")
		if req, ok := request.(*UpdateRequestInfo); ok {
			svcRequestInfo := service.NewUpdateRequestInfo(ctx, req.UpdateType, req.ServerName, req.Token, req.Data)
			svcResponseInfo := updater.Update(svcRequestInfo)
			updateResponseInfo := newUpdateResponseInfo(svcResponseInfo)
			return updateResponseInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

func MakeWatcherEndpoint(watcher service.Watcher) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//utils.Logger.Debug("request to bifrost watcher endpoint")
		if req, ok := request.(*WatchRequestInfo); ok {
			svcRequestInfo := service.NewWatchRequestInfo(ctx, req.WatchType, req.ServerName, req.Token, req.WatchObject)
			svcResponseInfo := watcher.Watch(svcRequestInfo)
			if svcResponseInfo.Error() != nil {
				return nil, svcResponseInfo.Error()
			}
			watchRespInfo := newWatchResponseInfo(svcResponseInfo, make(chan int), make(chan BytesResponseInfo))
			return watchRespInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

// HealthRequestInfo 健康检查请求结构
type HealthRequestInfo struct{}

// HealthResponseInfo 健康检查响应结构
type HealthResponseInfo struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if svc == nil {
			return HealthResponseInfo{false}, ErrInvalidService
		}
		return HealthResponseInfo{true}, nil
	}
}
