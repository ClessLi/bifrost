package endpoint

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

var (
	ErrInvalidRequest = errors.New("request has only one class: Request")
)

// ViewRequestInfo 信息展示请求信息
type ViewRequestInfo struct {
	ViewType   string `json:"view_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
}

// UpdateRequestInfo 数据更新请求信息
type UpdateRequestInfo struct {
	UpdateType string `json:"update_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
	Data       []byte `json:"data"`
}

// WatchRequestInfo 数据监看请求信息
type WatchRequestInfo struct {
	WatchType   string `json:"watch_type"`
	ServerName  string `json:"server_name"`
	Token       string `json:"token"`
	WatchObject string `json:"watch_object"`
}

// BytesResponseInfo 字节数据反馈信息接口对象
type BytesResponseInfo interface {
	Respond() []byte
	Error() string
}

// ErrorResponseInfo 错误反馈信息接口对象
type ErrorResponseInfo interface {
	Error() string
}

// WatchResponseInfo 数据监看反馈信息接口对象
type WatchResponseInfo interface {
	Respond() <-chan BytesResponseInfo
	Close() error
}

type bytesResponseInfo struct {
	Result *bytes.Buffer `json:"result"`
	Err    error         `json:"error"`
}

func (br bytesResponseInfo) Respond() []byte {
	return br.Result.Bytes()
}

func (br bytesResponseInfo) Error() string {
	if br.Err != nil {
		return br.Err.Error()
	}
	return ""
}

type errorResponseInfo struct {
	Err error `json:"error"`
}

func (er errorResponseInfo) Error() string {
	if er.Err != nil {
		return er.Err.Error()
	}
	return ""
}

type watchResponseInfo struct {
	Result    chan BytesResponseInfo `json:"result"`
	closeFunc func() error
}

func (wr watchResponseInfo) Respond() <-chan BytesResponseInfo {
	return wr.Result
}

func (wr watchResponseInfo) Close() error {
	return wr.closeFunc()
}

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
		utils.Logger.Debug("request to bifrost viewer endpoint")
		if req, ok := request.(*ViewRequestInfo); ok {
			svcRequestInfo := service.NewViewRequestInfo(ctx, req.ViewType, req.ServerName, req.Token)
			svcResponseInfo := viewer.View(svcRequestInfo)
			viewResponseInfo := newViewResponseInfo(svcResponseInfo)
			return viewResponseInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newViewResponseInfo(svcResponseInfo service.ViewResponseInfo) BytesResponseInfo {
	return &bytesResponseInfo{
		Result: bytes.NewBuffer(svcResponseInfo.Bytes()),
		Err:    svcResponseInfo.Error(),
	}
}

func MakeUpdaterEndpoint(updater service.Updater) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		utils.Logger.Debug("request to bifrost updater endpoint")
		if req, ok := request.(*UpdateRequestInfo); ok {
			svcRequestInfo := service.NewUpdateRequestInfo(ctx, req.UpdateType, req.ServerName, req.Token, req.Data)
			svcResponseInfo := updater.Update(svcRequestInfo)
			updateResponseInfo := newUpdateResponseInfo(svcResponseInfo)
			return updateResponseInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newUpdateResponseInfo(svcResponseInfo service.UpdateResponseInfo) ErrorResponseInfo {
	return &errorResponseInfo{Err: svcResponseInfo.Error()}
}

func MakeWatcherEndpoint(watcher service.Watcher) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		utils.Logger.Debug("request to bifrost watcher endpoint")
		if req, ok := request.(*WatchRequestInfo); ok {
			svcRequestInfo := service.NewWatchRequestInfo(ctx, req.WatchType, req.ServerName, req.Token, req.WatchObject)
			svcResponseInfo := watcher.Watch(svcRequestInfo)
			if svcResponseInfo.Error() != nil {
				return nil, svcResponseInfo.Error()
			}
			watchRespInfo := newWatchResponseInfo(svcResponseInfo)
			return watchRespInfo, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newWatchResponseInfo(svcResponseInfo service.WatchResponseInfo) WatchResponseInfo {
	signalChan := make(chan int)
	closeFunc := func() error {
		signalChan <- 9
		return nil
	}
	BytesResponseInfoChan := make(chan BytesResponseInfo)
	go func() {
		for {
			select {
			case b := <-svcResponseInfo.BytesChan():
				BytesResponseInfoChan <- &bytesResponseInfo{
					Result: bytes.NewBuffer(b),
					Err:    nil,
				}
			case err := <-svcResponseInfo.TransferErrorChan():
				BytesResponseInfoChan <- &bytesResponseInfo{
					Result: bytes.NewBuffer([]byte("")),
					Err:    err,
				}
			case sig := <-signalChan:
				if sig == 9 {
					err := svcResponseInfo.Close()
					if err != nil {
						utils.Logger.WarningF("[%s] service watch responseInfo close error, cased by %s", svcResponseInfo.GetServerName(), err)
					}
					goto stopHere
				}
			}
		}
	stopHere:
		return
	}()
	return &watchResponseInfo{
		Result:    BytesResponseInfoChan,
		closeFunc: closeFunc,
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{true}, nil
	}
}
