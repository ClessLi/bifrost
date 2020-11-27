package endpoint

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var (
	ErrInvalidVCReqType       = errors.New("RequestType has only one type: ViewConfig")
	ErrInvalidGCReqType       = errors.New("RequestType has only one type: GetConfig")
	ErrInvalidUCReqType       = errors.New("RequestType has only one type: UpdateConfig")
	ErrInvalidVSReqType       = errors.New("RequestType has only one type: ViewStatistics")
	ErrInvalidStatusReqType   = errors.New("RequestType has only one type: Status")
	ErrInvalidOperateRequest  = errors.New("request has only one class: OperateRequest")
	ErrInvalidConfigRequest   = errors.New("request has only one class: ConfigRequest")
	ErrInvalidWatchLogRequest = errors.New("request has only one class: WatchLogRequest")
	ErrResponseNull           = errors.New("response is null")
	ErrInvalidResponse        = errors.New("response is invalid")
)

// The service.Service method of BifrostEndpoints is used for the endpoint of the client
type BifrostEndpoints struct {
	ViewConfigEndpoint     endpoint.Endpoint
	GetConfigEndpoint      endpoint.Endpoint
	UpdateConfigEndpoint   endpoint.Endpoint
	ViewStatisticsEndpoint endpoint.Endpoint
	StatusEndpoint         endpoint.Endpoint
	WatchLogEndpoint       endpoint.Endpoint
	//StopWatchLogEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (ue BifrostEndpoints) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.ViewConfigEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: svrName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*bifrostpb.OperateResponse); ok {
		return response.Ret, err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostEndpoints) GetConfig(ctx context.Context, token, srvName string) (jsonData []byte, err error) {
	resp, err := ue.GetConfigEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: srvName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*bifrostpb.ConfigResponse); ok {
		return response.Ret.JData, err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostEndpoints) UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error) {
	resp, err := ue.UpdateConfigEndpoint(ctx, &bifrostpb.ConfigRequest{
		Token:   token,
		SvrName: svrName,
		Req: &bifrostpb.Config{
			JData: jsonData,
		},
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*bifrostpb.OperateResponse); ok {
		return response.Ret, nil
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostEndpoints) ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	resp, err := ue.ViewStatisticsEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: svrName,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*bifrostpb.OperateResponse); ok {
		return response.Ret, err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostEndpoints) Status(ctx context.Context, token string) (jsonData []byte, err error) {
	resp, err := ue.StatusEndpoint(ctx, &bifrostpb.OperateRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}
	if response, ok := resp.(*bifrostpb.OperateResponse); ok {
		return response.Ret, err
	} else {
		return nil, ErrResponseNull
	}
}

func (ue BifrostEndpoints) WatchLog(ctx context.Context, token, svrName, logName string, dataChan chan<- []byte, signal <-chan int) (err error) {
	// 客户端日志监看方法调用endpoint获取到gRPC客户端流对象
	res, err := ue.WatchLogEndpoint(ctx, nil)
	if err != nil {
		return err
	}
	// 判断endpoint传回并编码回的对象是否为gRPC客户端流对象
	if stream, ok := res.(bifrostpb.BifrostService_WatchLogClient); ok {
		//fmt.Println("初始化gRPC客户端成功")
		// 初始化gRPC服务返回管道和错误返回管道
		respChan := make(chan *bifrostpb.OperateResponse, 1)
		errChan := make(chan error, 1)
		// 创建接收匿名函数
		recv := func() {
			r, err := stream.Recv()
			if err != nil {
				errChan <- err
			}
			respChan <- r
		}
		// 创建客户端传出匿名函数
		sendOut := func(data []byte) error {
			select {
			case dataChan <- data:
				return nil
			case <-time.After(time.Second * 30): // 30秒传出未被接收将超时
				return service.ErrDataSendingTimeout
			}
		}
		// 发起request请求
		//fmt.Println("向gRPC服务端发送请求")
		err = stream.Send(&bifrostpb.OperateRequest{
			Token:   token,
			SvrName: svrName,
			Param:   logName,
		})
		if err != nil {
			return err
		}

		// 进入日志数据接收循环
		defer func() {
			err = stream.CloseSend()
		}()
		for {
			go recv() // 协程接收数据，包括数据和错误
			select {
			case sig := <-signal: // 客户端外边传入终止信号
				if sig == 9 {
					//fmt.Println("client shut down...")
					return err
				}
			case resp := <-respChan: // 接收到数据
				//fmt.Println("接收到数据")
				err = sendOut(resp.Ret) // 客户端日志数据传出
				if err != nil {
					return
				}
			case err = <-errChan: // 接收到gRPC服务错误返回
				return err
			}
		}
	}
	return ErrInvalidResponse
}

type OperateRequest struct {
	RequestType string `json:"request_type"`
	Token       string `json:"token"`
	SvrName     string `json:"svr_name"`
	Param       string `json:"param"`
}

type OperateResponse struct {
	Result []byte `json:"result"`
	Error  error  `json:"error"`
}

type WatchLogRequest struct {
	Token      string       `json:"token"`
	SvrName    string       `json:"svr_name"`
	Param      string       `json:"param"`
	DataChan   *chan []byte `json:"data_chan"`
	SignalChan *chan int    `json:"signal_chan"`
}

func NewWatchLogRequest(r *bifrostpb.OperateRequest) *WatchLogRequest {
	dataChan := make(chan []byte, 1)
	signalChan := make(chan int, 1)
	return &WatchLogRequest{
		Token:      r.Token,
		SvrName:    r.SvrName,
		Param:      r.Param,
		DataChan:   &dataChan,
		SignalChan: &signalChan,
	}
}

//type WatchLogResponse struct {
//	Result *chan []byte `json:"result"`
//	Signal *chan int `json:"signal"`
//	ErrChan *chan error `json:"err_chan"`
//}

type Config struct {
	JData []byte `json:"j_data"`
}

type ConfigRequest struct {
	RequestType string `json:"request_type"`
	Token       string `json:"token"`
	SvrName     string `json:"svr_name"`
	Ret         Config `json:"ret"`
}

type ConfigResponse struct {
	Result Config `json:"result"`
	Error  error  `json:"error"`
}

func MakeViewConfigEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(OperateRequest); ok {
			if strings.EqualFold(req.RequestType, "ViewConfig") {
				res, err := svc.ViewConfig(ctx, req.Token, req.SvrName)
				return OperateResponse{
					Result: res,
					Error:  err,
				}, nil
			}
			return nil, ErrInvalidVCReqType
		}
		return nil, ErrInvalidOperateRequest
	}
}

func MakeGetConfigEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(OperateRequest); ok {
			if strings.EqualFold(req.RequestType, "GetConfig") {
				res, err := svc.GetConfig(ctx, req.Token, req.SvrName)
				return ConfigResponse{
					Result: Config{JData: res},
					Error:  err,
				}, nil
			}
			return nil, ErrInvalidGCReqType
		}
		return nil, ErrInvalidOperateRequest
	}
}

func MakeUpdateConfigEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(ConfigRequest); ok {
			if strings.EqualFold(req.RequestType, "UpdateConfig") {
				res, err := svc.UpdateConfig(ctx, req.Token, req.SvrName, req.Ret.JData)
				return OperateResponse{
					Result: res,
					Error:  err,
				}, nil
			}
			return nil, ErrInvalidUCReqType
		}
		return nil, ErrInvalidConfigRequest
	}
}

func MakeViewStatisticsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(OperateRequest); ok {
			if strings.EqualFold(req.RequestType, "ViewStatistics") {
				res, err := svc.ViewStatistics(ctx, req.Token, req.SvrName)
				return OperateResponse{
					Result: res,
					Error:  err,
				}, nil
			}
			return nil, ErrInvalidVSReqType
		}
		return nil, ErrInvalidOperateRequest
	}
}

func MakeStatusEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(OperateRequest); ok {
			if strings.EqualFold(req.RequestType, "Status") {
				res, err := svc.Status(ctx, req.Token)
				return OperateResponse{
					Result: res,
					Error:  err,
				}, nil
			}
			return nil, ErrInvalidStatusReqType
		}
		return nil, ErrInvalidOperateRequest
	}
}

func MakeWatchLogEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//if req, ok := request.(OperateRequest); ok {
		//	if strings.EqualFold(req.RequestType, "WatchLog") {
		//		dataChan := make(chan []byte, 1)
		//		signal := make(chan int, 1)
		//		errChan := make(chan error, 1)
		//		// 重复进程异常
		//		go func() {
		//			errChan <- svc.WatchLog(ctx, req.Token, req.SvrName, req.Param, dataChan, signal)
		//		}()
		//		return WatchLogResponse{
		//			Result: &dataChan,
		//			Signal: &signal,
		//			ErrChan: &errChan,
		//		}, nil
		//	}
		//	return nil, ErrInvalidWLReqType
		//}
		//return nil, ErrInvalidOperateRequest
		fmt.Println("接受到handler请求")
		if req, ok := request.(*WatchLogRequest); ok {
			fmt.Println("请求发往service处理")
			err = svc.WatchLog(ctx, req.Token, req.SvrName, req.Param, *req.DataChan, *req.SignalChan)
			fmt.Println("service处理完毕")
			return OperateResponse{
				Result: nil,
				Error:  err,
			}, nil
		}
		fmt.Printf("%T\n", request)
		return nil, ErrInvalidWatchLogRequest
	}
}

func MakeWatchLogClientEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		client := bifrostpb.NewBifrostServiceClient(conn)
		return client.WatchLog(ctx)
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
