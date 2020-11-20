package endpoint

import (
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"strings"
)

var (
	ErrInvalidOperateRequestType = errors.New("RequestType has only four type: ViewConfig, GetConfig, ViewStatistics, Status")
	ErrInvalidConfigRequestType  = errors.New("RequestType has only one type: UpdateConfig")
	ErrInvalidRequest            = errors.New("request has only two class: OperateRequest, ConfigRequest")
)

type BifrostEndpoints struct {
	BifrostEndpoint     endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (ue BifrostEndpoints) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	resp, err := ue.BifrostEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: svrName,
	})

	response := resp.(*bifrostpb.OperateResponse)
	return response.Ret, err
}

func (ue BifrostEndpoints) GetConfig(ctx context.Context, token, srvName string) (jsonData []byte, err error) {
	resp, err := ue.BifrostEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: srvName,
	})

	response := resp.(*bifrostpb.ConfigResponse)
	return response.Ret.JData, err
}

func (ue BifrostEndpoints) UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error) {
	resp, err := ue.BifrostEndpoint(ctx, &bifrostpb.ConfigRequest{
		Token:   token,
		SvrName: svrName,
		Req: &bifrostpb.Config{
			JData: jsonData,
		},
	})
	response := resp.(*bifrostpb.OperateResponse)
	return response.Ret, err
}

func (ue BifrostEndpoints) ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	resp, err := ue.BifrostEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: svrName,
	})

	response := resp.(*bifrostpb.OperateResponse)
	return response.Ret, err
}

func (ue BifrostEndpoints) Status(ctx context.Context, token string) (jsonData []byte, err error) {
	resp, err := ue.BifrostEndpoint(ctx, &bifrostpb.OperateRequest{
		Token:   token,
		SvrName: "",
	})

	response := resp.(*bifrostpb.OperateResponse)
	return response.Ret, err
}

type OperateRequest struct {
	RequestType string `json:"request_type"`
	Token       string `json:"token"`
	SvrName     string `json:"svr_name"`
}

type OperateResponse struct {
	Result []byte `json:"result"`
	Error  error  `json:"error"`
}

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

func MakeBifrostEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var res []byte
		switch request.(type) {
		case OperateRequest:
			req := request.(OperateRequest)
			if strings.EqualFold(req.RequestType, "ViewConfig") {
				res, err = svc.ViewConfig(ctx, req.Token, req.SvrName)
			} else if strings.EqualFold(req.RequestType, "GetConfig") {
				res, err = svc.GetConfig(ctx, req.Token, req.SvrName)
				return ConfigResponse{
					Result: Config{JData: res},
					Error:  err,
				}, nil
			} else if strings.EqualFold(req.RequestType, "ViewStatistics") {
				res, err = svc.ViewStatistics(ctx, req.Token, req.SvrName)
			} else if strings.EqualFold(req.RequestType, "Status") {
				res, err = svc.Status(ctx, req.Token)
			} else {
				return nil, ErrInvalidOperateRequestType
			}

			return OperateResponse{
				Result: res,
				Error:  err,
			}, nil
		case ConfigRequest:
			req := request.(ConfigRequest)
			if strings.EqualFold(req.RequestType, "UpdateConfig") {
				res, err = svc.UpdateConfig(ctx, req.Token, req.SvrName, req.Ret.JData)
			} else {
				return nil, ErrInvalidConfigRequestType
			}
			return OperateResponse{
				Result: res,
				Error:  err,
			}, nil
		default:
			return nil, ErrInvalidRequest
		}
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
