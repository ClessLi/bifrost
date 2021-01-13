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
	ErrInvalidRequest  = errors.New("request has only one class: Request")
	ErrResponseNull    = errors.New("response is null")
	ErrInvalidResponse = errors.New("response is invalid")
	ErrUnknownResponse = errors.New("unknown response")
)

type ViewRequest struct {
	ViewType   string `json:"view_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
}

type UpdateRequest struct {
	UpdateType string `json:"update_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
	Data       []byte `json:"data"`
}

type WatchRequest struct {
	WatchType   string `json:"watch_type"`
	ServerName  string `json:"server_name"`
	Token       string `json:"token"`
	WatchObject string `json:"watch_object"`
}

type BytesResponder interface {
	Respond() []byte
	Error() string
}

type ErrorResponder interface {
	Error() string
}

type WatchResponder interface {
	Respond() <-chan BytesResponder
	Close() error
}

type bytesResponder struct {
	Result *bytes.Buffer `json:"result"`
	Err    error         `json:"error"`
}

func (br bytesResponder) Respond() []byte {
	return br.Result.Bytes()
}

func (br bytesResponder) Error() string {
	if br.Err != nil {
		return br.Err.Error()
	}
	return ""
}

type errorResponder struct {
	Err error `json:"error"`
}

func (er errorResponder) Error() string {
	if er.Err != nil {
		return er.Err.Error()
	}
	return ""
}

type watchResponder struct {
	Result    chan BytesResponder `json:"result"`
	closeFunc func() error
}

func (wr watchResponder) Respond() <-chan BytesResponder {
	return wr.Result
}

func (wr watchResponder) Close() error {
	return wr.closeFunc()
}

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
		if req, ok := request.(*ViewRequest); ok {
			svcRequester := service.NewViewRequester(ctx, req.ViewType, req.ServerName, req.Token)
			svcResponder := viewer.View(svcRequester)
			viewResponder := newViewResponder(svcResponder)
			return viewResponder, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newViewResponder(svcResponder service.ViewResponder) BytesResponder {
	return &bytesResponder{
		Result: bytes.NewBuffer(svcResponder.Bytes()),
		Err:    svcResponder.Error(),
	}
}

func MakeUpdaterEndpoint(updater service.Updater) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		utils.Logger.Debug("request to bifrost updater endpoint")
		if req, ok := request.(*UpdateRequest); ok {
			svcRequester := service.NewUpdateRequester(ctx, req.UpdateType, req.ServerName, req.Token, req.Data)
			svcResponder := updater.Update(svcRequester)
			updateResponder := newUpdateResponder(svcResponder)
			return updateResponder, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newUpdateResponder(svcResponder service.UpdateResponder) ErrorResponder {
	return &errorResponder{Err: svcResponder.Error()}
}

func MakeWatcherEndpoint(watcher service.Watcher) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		utils.Logger.Debug("request to bifrost watcher endpoint")
		if req, ok := request.(*WatchRequest); ok {
			svcRequester := service.NewWatchRequester(ctx, req.WatchType, req.ServerName, req.Token, req.WatchObject)
			svcResponder := watcher.Watch(svcRequester)
			if svcResponder.Error() != nil {
				return nil, svcResponder.Error()
			}
			watchResponder := newWatchResponder(svcResponder)
			return watchResponder, nil
		}
		return nil, ErrInvalidRequest
	}
}

func newWatchResponder(svcResponder service.WatchResponder) WatchResponder {
	signalChan := make(chan int)
	closeFunc := func() error {
		signalChan <- 9
		return nil
	}
	BytesResponderChan := make(chan BytesResponder)
	go func() {
		for {
			select {
			case b := <-svcResponder.BytesChan():
				BytesResponderChan <- &bytesResponder{
					Result: bytes.NewBuffer(b),
					Err:    nil,
				}
			case err := <-svcResponder.TransferErrorChan():
				BytesResponderChan <- &bytesResponder{
					Result: bytes.NewBuffer([]byte("")),
					Err:    err,
				}
			case sig := <-signalChan:
				if sig == 9 {
					err := svcResponder.Close()
					utils.Logger.WarningF("[%s] service watch responder close error, cased by %s", svcResponder.GetServerName(), err)
					goto stopHere
				}
			}
		}
	stopHere:
		return
	}()
	return &watchResponder{
		Result:    BytesResponderChan,
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
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{true}, nil
	}
}
