package endpoint

import (
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

// The service.Service method of BifrostEndpoints is used for the endpoint of the client
type BifrostClientEndpoints struct {
	getInfo   endpoint.Endpoint
	update    endpoint.Endpoint
	watchInfo endpoint.Endpoint
}

func NewBifrostClient(conn *grpc.ClientConn) *BifrostClientEndpoints {
	return &BifrostClientEndpoints{
		getInfo: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"GetInfo",
			EncodeBifrostServiceClientRequest,
			DecodeBifrostServiceClientResponse,
			new(bifrostpb.Response),
		).Endpoint(),
		update: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"UpdateConfig",
			EncodeBifrostServiceClientRequest,
			DecodeBifrostServiceClientResponse,
			new(bifrostpb.Response),
		).Endpoint(),
		watchInfo: MakeWatchLogClientEndpoint(conn),
	}
}

func EncodeBifrostServiceClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return encodeClientRequest(ctx, r)
}

func DecodeBifrostServiceClientResponse(_ context.Context, r interface{}) (interface{}, error) {
	if resp, ok := r.(*bifrostpb.Response); ok {
		var respErr error
		if resp.Err != "" {
			respErr = errors.New(resp.Err)
		}
		return Response{
			Result: resp.Ret,
			Error:  respErr,
		}, nil
	}
	return nil, ErrUnknownResponse
}

func encodeClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(Request); ok {
		return &bifrostpb.Request{
			Token:       req.Token,
			SvrName:     req.ServerName,
			RequestType: req.RequestType,
			Param:       req.Param,
			Data:        req.Data,
		}, nil
	}
	return nil, ErrInvalidBifrostServiceEndpointRequest
}

func (ue BifrostClientEndpoints) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.getInfo(ctx, Request{
		RequestType: "DisplayConfig",
		Token:       token,
		ServerName:  svrName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(Response); ok {
		return response.Result, response.Error
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) GetConfig(ctx context.Context, token, srvName string) (data []byte, err error) {
	resp, err := ue.getInfo(ctx, Request{
		RequestType: "GetConfig",
		Token:       token,
		ServerName:  srvName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(Response); ok {
		return response.Result, response.Error
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) UpdateConfig(ctx context.Context, token, svrName string, reqData []byte, params ...string) (data []byte, err error) {
	req := Request{
		RequestType: "UpdateConfig",
		Token:       token,
		ServerName:  svrName,
		Data:        reqData,
	}
	if params != nil {
		req.Param = params[0]
	} else {
		req.Param = "full"
	}
	resp, err := ue.update(ctx, req)

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(Response); ok {
		return response.Result, response.Error
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) ViewStatistics(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.getInfo(ctx, Request{
		RequestType: "ShowStatistics",
		Token:       token,
		ServerName:  svrName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(Response); ok {
		return response.Result, response.Error
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) Status(ctx context.Context, token string) (data []byte, err error) {
	resp, err := ue.getInfo(ctx, Request{
		RequestType: "DisplayStatus",
		Token:       token,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(Response); ok {
		return response.Result, response.Error
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostClientEndpoints) WatchLog(ctx context.Context, token, svrName, logName string) (logWatcher *Watcher, err error) {
	// 客户端日志监看方法调用endpoint获取到gRPC客户端流对象
	res, err := ue.watchInfo(ctx, nil)
	if err != nil {
		return nil, err
	}
	// 判断endpoint传回并编码回的对象是否为gRPC客户端流对象
	if stream, ok := res.(bifrostpb.BifrostService_WatchInfoClient); ok {
		//fmt.Println("初始化gRPC客户端成功")
		// 初始化gRPC服务返回管道和错误返回管道
		respChan := make(chan *bifrostpb.Response)
		dataChan := make(chan []byte)
		errChan := make(chan error)
		sigChan := make(chan int)
		closeFunc := func() error {
			sigChan <- 9
			return nil
		}
		logWatcher = NewWatcher(dataChan, errChan, closeFunc)
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
		err = stream.Send(&bifrostpb.Request{
			RequestType: "WatchLog",
			Token:       token,
			SvrName:     svrName,
			Param:       logName,
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
		client := bifrostpb.NewBifrostServiceClient(conn)
		return client.WatchInfo(ctx)
	}
}
