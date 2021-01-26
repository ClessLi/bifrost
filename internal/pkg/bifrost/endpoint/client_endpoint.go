package endpoint

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var (
	ErrInvalidBifrostServiceEndpointRequest = errors.New("request has only one: endpoint.Request")
)

// The service.ServiceConfig method of BifrostEndpoints is used for the endpoint of the client
type BifrostClientEndpoints struct {
	viewEndpoint   endpoint.Endpoint
	updateEndpoint endpoint.Endpoint
	watchEndpoint  endpoint.Endpoint
}

func NewBifrostClient(conn *grpc.ClientConn) *BifrostClientEndpoints {
	return &BifrostClientEndpoints{
		viewEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.ViewService",
			"View",
			encodeClientRequest,
			decodeClientResponse,
			new(bifrostpb.BytesResponse),
		).Endpoint(),
		updateEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.UpdateService",
			"Update",
			encodeClientRequest,
			decodeClientResponse,
			new(bifrostpb.ErrorResponse),
		).Endpoint(),
		watchEndpoint: MakeWatchLogClientEndpoint(conn),
	}
}

//func EncodeViewServiceClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return encodeClientRequest(ctx, r)
//}
//
//func DecodeViewServiceClientResponse(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeClientResponse(ctx, r)
//}
//
//func EncodeUpdateServiceClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return encodeClientRequest(ctx, r)
//}

func encodeClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r.(type) {
	case ViewRequest:
		req := r.(ViewRequest)
		return &bifrostpb.ViewRequest{
			ViewType:   req.ViewType,
			ServerName: req.ServerName,
			Token:      req.Token,
		}, nil
	case UpdateRequest:
		req := r.(UpdateRequest)
		return &bifrostpb.UpdateRequest{
			UpdateType: req.UpdateType,
			ServerName: req.ServerName,
			Token:      req.Token,
			Data:       req.Data,
		}, nil
	}

	return nil, ErrInvalidBifrostServiceEndpointRequest
}

func decodeClientResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r.(type) {
	case *bifrostpb.BytesResponse:
		resp := r.(*bifrostpb.BytesResponse)
		var respErr error
		if resp.Err != "" {
			respErr = errors.New(resp.Err)
		}
		return &bytesResponder{
			Result: bytes.NewBuffer(resp.Ret),
			Err:    respErr,
		}, nil
	case *bifrostpb.ErrorResponse:
		resp := r.(*bifrostpb.ErrorResponse)
		var respErr error
		if resp.Err != "" {
			respErr = errors.New(resp.Err)
		}
		return &errorResponder{Err: respErr}, nil
	}
	return nil, ErrUnknownResponse
}

func (ue BifrostClientEndpoints) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.viewEndpoint(ctx, ViewRequest{
		ViewType:   "DisplayConfig",
		Token:      token,
		ServerName: svrName,
	})

	if err != nil {
		return nil, err
	}
	if responder, ok := resp.(*bytesResponder); ok {
		return responder.Respond(), responder.Err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) GetConfig(ctx context.Context, token, srvName string) (data []byte, err error) {
	resp, err := ue.viewEndpoint(ctx, ViewRequest{
		ViewType:   "GetConfig",
		Token:      token,
		ServerName: srvName,
	})

	if err != nil {
		return nil, err
	}
	if responder, ok := resp.(*bytesResponder); ok {
		return responder.Respond(), responder.Err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) UpdateConfig(ctx context.Context, token, svrName string, reqData []byte, params ...string) (data []byte, err error) {
	req := UpdateRequest{
		UpdateType: "UpdateConfig",
		Token:      token,
		ServerName: svrName,
		Data:       reqData,
	}
	resp, err := ue.updateEndpoint(ctx, req)

	if err != nil {
		return nil, err
	}
	if responder, ok := resp.(*errorResponder); ok {
		return nil, responder.Err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) ViewStatistics(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.viewEndpoint(ctx, ViewRequest{
		ViewType:   "ShowStatistics",
		Token:      token,
		ServerName: svrName,
	})

	if err != nil {
		return nil, err
	}
	if responder, ok := resp.(*bytesResponder); ok {
		return responder.Respond(), responder.Err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) Status(ctx context.Context, token string) (data []byte, err error) {
	resp, err := ue.viewEndpoint(ctx, ViewRequest{
		ViewType: "DisplayStatus",
		Token:    token,
	})

	if err != nil {
		return nil, err
	}
	if responder, ok := resp.(*bytesResponder); ok {
		return responder.Respond(), responder.Err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) WatchLog(ctx context.Context, token, svrName, logName string) (logWatcher ClientWatcher, err error) {
	// 客户端日志监看方法调用endpoint获取到gRPC客户端流对象
	res, err := ue.watchEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	// 判断endpoint传回并编码回的对象是否为gRPC客户端流对象
	if stream, ok := res.(bifrostpb.WatchService_WatchClient); ok {
		//fmt.Println("初始化gRPC客户端成功")
		// 初始化gRPC服务返回管道和错误返回管道
		respChan := make(chan *bifrostpb.BytesResponse)
		dataChan := make(chan []byte)
		errChan := make(chan error)
		sigChan := make(chan int)
		closeFunc := func() error {
			sigChan <- 9
			return nil
		}
		logWatcher = newClientLogWatcher(dataChan, errChan, closeFunc)
		// 创建接收匿名函数
		recv := func() {
			r, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			respChan <- r
		}
		// 创建客户端传出匿名函数
		sendOut := func(data []byte) error {
			select {
			case dataChan <- data:
				return nil
			case <-time.After(time.Second * 30): // 30秒传出未被接收将超时
				return web_server_manager.ErrDataSendingTimeout
			}
		}
		// 发起request请求
		//fmt.Println("向gRPC服务端发送请求")
		err = stream.Send(&bifrostpb.WatchRequest{
			WatchType:   "WatchLog",
			ServerName:  svrName,
			Token:       token,
			WatchObject: logName,
		})
		if err != nil {
			return nil, err
		}

		// 进入日志数据接收循环
		go func() {
			defer func() {
				err := stream.CloseSend()
				if err != nil {
					errChan <- err
				}
			}()
			for {
				go recv() // 协程接收数据，包括数据和错误
				select {
				case sig := <-sigChan: // 客户端外边传入终止信号
					if sig == 9 {
						//fmt.Println("client shut down...")
						return
					}
				case resp := <-respChan: // 接收到数据
					//fmt.Println("接收到数据")
					if resp.Ret != nil {
						err = sendOut(resp.Ret) // 客户端日志数据传出
						if err != nil {
							errChan <- err
							return
						}
					} else if resp.Err != "" {
						errChan <- errors.New(resp.Err)
						return
					} else {
						errChan <- errors.New("response is null")
					}
				}
			}
		}()
		return logWatcher, nil
	}
	return nil, ErrInvalidResponse
}

func MakeWatchLogClientEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		client := bifrostpb.NewWatchServiceClient(conn)
		return client.Watch(ctx)
	}
}
