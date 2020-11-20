package transport

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"io"
	"time"
)

var (
	ChunkSize = 1024

	recvTimeout time.Duration = 5 // 5 minutes

	ErrRequestInconsistent = errors.New("the request is inconsistent")
	ErrRecvTimeout         = errors.New("receive timeout during data wrap")
	ErrUnknownRequest      = errors.New("unknown request error")
	//ErrUnknownResponse       = errors.New("unknown response error")
	ErrInvalidOperateRequest = errors.New("request has only one: OperateRequest")
	ErrInvalidConfigRequest  = errors.New("request has only one: ConfigRequest")
)

type grpcServer struct {
	viewConfig     grpc.Handler
	getConfig      grpc.Handler
	updateConfig   grpc.Handler
	viewStatistics grpc.Handler
	status         grpc.Handler
}

func (s *grpcServer) ViewConfig(r *bifrostpb.OperateRequest, stream bifrostpb.BifrostService_ViewConfigServer) error {
	_, resp, err := s.viewConfig.ServeGRPC(stream.Context(), r)
	if err != nil {
		return stream.Send(&bifrostpb.OperateResponse{
			Ret: nil,
			Err: err.Error(),
		})
	}
	response := resp.(*bifrostpb.OperateResponse)
	n := len(response.Ret)
	for i := 0; i < n; i += ChunkSize {
		// TODO: response优化为[][]byte，用BifrostService.ChunckSize定义每个[]byte切片容积
		if n <= i+ChunkSize {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i:],
				Err: "",
			})
		} else {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i : i+ChunkSize],
				Err: "",
			})
		}
		if err != nil {
			return stream.Send(&bifrostpb.OperateResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}
	return nil
}

func (s *grpcServer) GetConfig(r *bifrostpb.OperateRequest, stream bifrostpb.BifrostService_GetConfigServer) error {
	_, resp, err := s.getConfig.ServeGRPC(stream.Context(), r)
	if err != nil {
		return stream.Send(&bifrostpb.ConfigResponse{
			Ret: nil,
			Err: err.Error(),
		})
	}
	response := resp.(*bifrostpb.ConfigResponse)
	n := len(response.Ret.JData)
	for i := 0; i < n; i += ChunkSize {
		if n <= i+ChunkSize {
			err = stream.Send(&bifrostpb.ConfigResponse{
				Ret: &bifrostpb.Config{JData: response.Ret.JData[i:]},
				Err: "",
			})
		} else {
			err = stream.Send(&bifrostpb.ConfigResponse{
				Ret: &bifrostpb.Config{JData: response.Ret.JData[i : i+ChunkSize]},
				Err: "",
			})
		}
		if err != nil {
			return stream.Send(&bifrostpb.ConfigResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}
	return nil
}

func (s *grpcServer) UpdateConfig(stream bifrostpb.BifrostService_UpdateConfigServer) error {
	var err error
	defer func() {
		if err != nil {
			_ = stream.SendAndClose(&bifrostpb.OperateResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}()
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	defer buffer.Reset()
	var req *bifrostpb.ConfigRequest
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
			req = &bifrostpb.ConfigRequest{
				Token:   in.Token,
				SvrName: in.SvrName,
				Req:     &bifrostpb.Config{},
			}
		} else {
			if in.SvrName != req.SvrName || in.Token != req.Token {
				err = ErrRequestInconsistent
				return err
			}
		}
		buffer.Write(in.Req.JData)
	}
	if isTimeout {
		err = ErrRecvTimeout
		return err
	}
	if req == nil {
		err = ErrUnknownRequest
		return err
	}
	req.Req.JData = buffer.Bytes()
	_, resp, err := s.updateConfig.ServeGRPC(stream.Context(), req)
	if err != nil {
		return stream.SendAndClose(&bifrostpb.OperateResponse{Err: err.Error()})
	}
	response := resp.(*bifrostpb.OperateResponse)
	return stream.SendAndClose(response)
}

func (s *grpcServer) ViewStatistics(r *bifrostpb.OperateRequest, stream bifrostpb.BifrostService_ViewStatisticsServer) error {
	_, resp, err := s.viewStatistics.ServeGRPC(stream.Context(), r)
	if err != nil {
		return stream.Send(&bifrostpb.OperateResponse{
			Ret: nil,
			Err: err.Error(),
		})
	}
	response := resp.(*bifrostpb.OperateResponse)
	n := len(response.Ret)
	for i := 0; i < n; i += ChunkSize {
		if n <= i+ChunkSize {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i:],
				Err: "",
			})
		} else {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i : i+ChunkSize],
				Err: "",
			})
		}
		if err != nil {
			return stream.Send(&bifrostpb.OperateResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}
	return nil
}

func (s *grpcServer) Status(r *bifrostpb.OperateRequest, stream bifrostpb.BifrostService_StatusServer) error {
	_, resp, err := s.status.ServeGRPC(stream.Context(), r)
	if err != nil {
		return stream.Send(&bifrostpb.OperateResponse{
			Ret: nil,
			Err: err.Error(),
		})
	}
	response := resp.(*bifrostpb.OperateResponse)
	n := len(response.Ret)
	for i := 0; i < n; i += ChunkSize {
		if n <= i+ChunkSize {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i:],
				Err: "",
			})
		} else {
			err = stream.Send(&bifrostpb.OperateResponse{
				Ret: response.Ret[i : i+ChunkSize],
				Err: "",
			})
		}
		if err != nil {
			return stream.Send(&bifrostpb.OperateResponse{
				Ret: nil,
				Err: err.Error(),
			})
		}
	}
	return nil
}

func NewBifrostServer(ctx context.Context, endpoints endpoint.BifrostEndpoints) bifrostpb.BifrostServiceServer {
	return &grpcServer{
		viewConfig:     grpc.NewServer(endpoints.ViewConfigEndpoint, DecodeViewConfigRequest, EncodeOperateResponse),
		getConfig:      grpc.NewServer(endpoints.GetConfigEndpoint, DecodeGetConfigRequest, EncodeConfigResponse), // dong大的坑，o(╥﹏╥)o，handler得绑定对应endpoint
		updateConfig:   grpc.NewServer(endpoints.UpdateConfigEndpoint, DecodeUpdateConfigRequest, EncodeOperateResponse),
		viewStatistics: grpc.NewServer(endpoints.ViewStatisticsEndpoint, DecodeViewStatisticsRequest, EncodeOperateResponse),
		status:         grpc.NewServer(endpoints.StatusEndpoint, DecodeStatusRequest, EncodeOperateResponse),
	}
}
