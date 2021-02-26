package auth

import (
	"context"
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/internal/pkg/auth/endpoint"
	"github.com/ClessLi/bifrost/internal/pkg/auth/service"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	//*grpc.ClientConn
	service.Service
}

func NewClient(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, err
	}

	eps := endpoint.AuthEndpoints{
		LoginEndpoint: grpctransport.NewClient(conn,
			"authpb.AuthService",
			"Login",
			decodeRequest,
			encodeResponse,
			authpb.AuthResponse{}, // 写成request会导致panic，is *authpb.AuthResponse, not *authpb.AuthRequest
		).Endpoint(),
		VerifyEndpoint: grpctransport.NewClient(conn,
			"authpb.AuthService",
			"Verify",
			decodeRequest,
			encodeResponse,
			authpb.VerifyResponse{}, // 写成request会导致panic，is *authpb.VerifyResponse, not *authpb.VerifyRequest
		).Endpoint(),
	}
	return &Client{
		Service: eps,
	}, nil
}

func decodeRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
	return r, nil
}

func encodeResponse(ctx context.Context, r interface{}) (response interface{}, err error) {
	return r, nil
}
