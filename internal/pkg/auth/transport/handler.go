package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/internal/pkg/auth/endpoint"

	"github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

type grpcServer struct {
	login  grpc.Handler
	verify grpc.Handler
}

func (s *grpcServer) Login(ctx context.Context, r *authpb.AuthRequest) (*authpb.AuthResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, r)
	if resp != nil {
		return resp.(*authpb.AuthResponse), err
	}
	//if err != nil {
	//	return nil, err
	//}
	return nil, err
}

func (s *grpcServer) Verify(ctx context.Context, r *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	_, resp, err := s.verify.ServeGRPC(ctx, r)
	//if err != nil {
	//	return nil, err
	//}
	//return resp.(*authpb.VerifyResponse), nil
	if resp != nil {
		return resp.(*authpb.VerifyResponse), err
	}
	return nil, err
}

func NewAuthServer(ctx context.Context, endpoints endpoint.AuthEndpoints) authpb.AuthServiceServer {
	return &grpcServer{
		login: grpc.NewServer(
			endpoints.LoginEndpoint,
			DecodeAuthRequest,
			EncodeAuthResponse,
		),
		verify: grpc.NewServer(
			endpoints.VerifyEndpoint,
			DecodeVerifyRequest,
			EncodeVerifyResponse,
		),
	}
}
