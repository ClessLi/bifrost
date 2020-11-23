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
	ErrInvalidVCReqType      = errors.New("RequestType has only one type: ViewConfig")
	ErrInvalidGCReqType      = errors.New("RequestType has only one type: GetConfig")
	ErrInvalidUCReqType      = errors.New("RequestType has only one type: UpdateConfig")
	ErrInvalidVSReqType      = errors.New("RequestType has only one type: ViewStatistics")
	ErrInvalidStatusReqType  = errors.New("RequestType has only one type: Status")
	ErrInvalidOperateRequest = errors.New("request has only one class: OperateRequest")
	ErrInvalidConfigRequest  = errors.New("request has only one class: ConfigRequest")
	ErrResponseNull          = errors.New("response is null")
)

type BifrostEndpoints struct {
	ViewConfigEndpoint     endpoint.Endpoint
	GetConfigEndpoint      endpoint.Endpoint
	UpdateConfigEndpoint   endpoint.Endpoint
	ViewStatisticsEndpoint endpoint.Endpoint
	StatusEndpoint         endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
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
		Token:   token,
		SvrName: "",
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
