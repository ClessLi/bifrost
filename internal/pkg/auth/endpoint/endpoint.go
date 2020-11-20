package endpoint

import (
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/internal/pkg/auth/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"strings"
)

type AuthEndpoints struct {
	AuthEndpoint        endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (ue AuthEndpoints) Login(ctx context.Context, username, password string, unexpired bool) (string, error) {
	resp, err := ue.AuthEndpoint(ctx, &authpb.AuthRequest{
		Username:  username,
		Password:  password,
		Unexpired: unexpired,
	})
	response := resp.(*authpb.AuthResponse)
	return response.Token, err
}

func (ue AuthEndpoints) Verify(ctx context.Context, token string) (bool, error) {
	resp, err := ue.AuthEndpoint(ctx, &authpb.VerifyRequest{Token: token})
	response := resp.(*authpb.VerifyResponse)
	return response.Passed, err
}

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Login, Verify")
	ErrInvalidRequest     = errors.New("request has only two class: AuthRequest, VerifyRequest")
)

type AuthRequest struct {
	RequestType string `json:"request_type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Unexpired   bool   `json:"unexpired"`
}

type AuthResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

type VerifyRequest struct {
	ResquesType string `json:"resques_type"`
	Token       string `json:"token"`
}

type VerifyResponse struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}

func MakeAuthEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		switch request.(type) {
		case AuthRequest:
			req := request.(AuthRequest)
			var res string
			if strings.EqualFold(req.RequestType, "Login") {
				res, err = svc.Login(ctx, req.Username, req.Password, req.Unexpired)
			} else {
				return nil, ErrInvalidRequestType
			}
			return AuthResponse{
				Result: res,
				Error:  err,
			}, nil

		case VerifyRequest:
			req := request.(VerifyRequest)
			var res bool
			if strings.EqualFold(req.ResquesType, "Verify") {
				res, _ = svc.Verify(ctx, req.Token)
			} else {
				return nil, ErrInvalidRequestType
			}
			return VerifyResponse{
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
