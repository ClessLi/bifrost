package transport

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"time"
)

var (
	ChunkSize = 1024

	recvTimeout time.Duration = 5 // 5 minutes

	ErrRequestInconsistent       = errors.New("the request is inconsistent")
	ErrRecvTimeout               = errors.New("receive timeout during data wrap")
	ErrUnknownRequest            = errors.New("unknown request error")
	ErrInvalidHealthCheckRequest = errors.New("request has only one: HealthCheckRequest")
	ErrUnknownResponse           = errors.New("unknown response error")

	ErrWatcherCloseTimeout = errors.New("close watcher timeout")
)

type gRPCHandlers struct {
	viewService   grpc.Handler
	updateService grpc.Handler
	watchService  grpc.Handler
}

type gRPCHealthCheckHandler struct {
	healthCheck grpc.Handler
}

func (s *gRPCHandlers) View(r *bifrostpb.ViewRequest, stream bifrostpb.ViewService_ViewServer) error {
	_, resp, err := s.viewService.ServeGRPC(stream.Context(), r)
	if err != nil {
		return stream.Send(&bifrostpb.BytesResponse{
			Ret: nil,
			Err: err.Error(),
		})
	}
	response := resp.(*bifrostpb.BytesResponse)
	n := len(response.Ret)
	for i := 0; i < n; i += ChunkSize {
		// TODO: response优化为[][]byte，用BifrostService.ChunckSize定义每个[]byte切片容积
		if n <= i+ChunkSize {
			err = stream.Send(&bifrostpb.BytesResponse{
				Ret: response.Ret[i:],
				Err: "",
			})
		} else {
			err = stream.Send(&bifrostpb.BytesResponse{
				Ret: response.Ret[i : i+ChunkSize],
				Err: "",
			})
		}
		if err != nil {
			return stream.Send(&bifrostpb.BytesResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}
	return nil
}

func (s *gRPCHandlers) Update(stream bifrostpb.UpdateService_UpdateServer) error {
	var err error
	defer func() {
		if err != nil {
			_ = stream.SendAndClose(&bifrostpb.ErrorResponse{
				Err: err.Error(),
			})
		}
	}()
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	defer buffer.Reset()
	var req *bifrostpb.UpdateRequest
	recvStartTime := time.Now()
	//recvTOPoint := time.Now().Unix() + recvTimeout
	var isTimeout bool
	for !isTimeout {
		//isTimeout = time.Now().Unix() < recvTOPoint
		isTimeout = time.Since(recvStartTime) >= recvTimeout*time.Minute
		in, err := stream.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		if req == nil {
			req = &bifrostpb.UpdateRequest{
				ServerName: in.ServerName,
				Token:      in.Token,
				UpdateType: in.UpdateType,
			}
		} else {
			if in.ServerName != req.ServerName || in.Token != req.Token || in.UpdateType != req.UpdateType {
				err = ErrRequestInconsistent
				return err
			}
		}
		buffer.Write(in.Data)
	}
	if isTimeout {
		err = ErrRecvTimeout
		return err
	}
	if req == nil {
		err = ErrUnknownRequest
		return err
	}
	req.Data = buffer.Bytes()
	_, resp, err := s.updateService.ServeGRPC(stream.Context(), req)
	if err != nil {
		return stream.SendAndClose(&bifrostpb.ErrorResponse{Err: err.Error()})
	}
	response := resp.(*bifrostpb.ErrorResponse)
	return stream.SendAndClose(response)
}

func (s *gRPCHandlers) Watch(stream bifrostpb.WatchService_WatchServer) (err error) {
	stopSendSig := make(chan int, 1)
	// 获取gRPC客户端请求
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	// 向endpoint发起请求，获取WatchLogResponse
	//fmt.Println("请求发往endpoint处理")
	_, resp, err := s.watchService.ServeGRPC(stream.Context(), req)
	if err != nil {
		return err
	}
	responder, ok := resp.(*watchResponder)
	if !ok {
		return ErrUnknownResponse
	}
	//fmt.Println("endpoint处理完毕")
	// 监听WatchLogResponse中Result(dataChan)和ErrChan，监听stopSendSig、WatchLogTimeout
	go func(stopSig chan int) {
		for {
			select {
			case response := <-responder.Respond():
				_ = stream.Send(response)
			case sig := <-stopSig:
				if sig == 9 { // 信号9传入则开始停止
					//fmt.Println("开始停止")
					//fmt.Println("开始发送停止信号")
					err = responder.Close()
					//fmt.Println("停止传递结束")
					return
				}
			}
		}
	}(stopSendSig)

	for {
		// 接收客户端请求
		in, err := stream.Recv()
		//fmt.Println("再次接受到客户端请求")
		if err == nil && (in.ServerName != req.ServerName || in.Token != req.Token || in.WatchType != req.WatchType || in.WatchObject != req.WatchObject) {
			err = ErrRequestInconsistent
		}
		if err != nil {
			stopSendSig <- 9
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}

func (s *gRPCHealthCheckHandler) Check(ctx context.Context, r *grpc_health_v1.HealthCheckRequest) (response *grpc_health_v1.HealthCheckResponse, err error) {
	_, resp, err := s.healthCheck.ServeGRPC(ctx, r)
	if err != nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, err
	}
	var ok bool
	if response, ok = resp.(*grpc_health_v1.HealthCheckResponse); ok {
		return response, nil
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
	}, nil
}

func (s *gRPCHealthCheckHandler) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func NewGRPCHandlers(ctx context.Context, endpoints endpoint.BifrostEndpoints) *gRPCHandlers {
	return &gRPCHandlers{
		viewService:   grpc.NewServer(endpoints.ViewerEndpoint, DecodeViewRequest, EncodeViewResponse),
		updateService: grpc.NewServer(endpoints.UpdaterEndpoint, DecodeUpdateRequest, EncodeUpdateResponse),
		watchService:  grpc.NewServer(endpoints.WatcherEndpoint, DecodeWatchRequest, EncodeWatchResponse),
	}
}

func NewHealthCheckHandler(ctx context.Context, endpoints endpoint.BifrostEndpoints) grpc_health_v1.HealthServer {
	return &gRPCHealthCheckHandler{
		healthCheck: grpc.NewServer(endpoints.HealthCheckEndpoint, DecodeHealthCheckRequest, EncodeHealthCheckResponse),
	}
}
