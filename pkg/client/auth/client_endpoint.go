package auth

import (
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrResponseNull = errors.New("response is null")
)

func decodeRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
	return r, nil
}

func encodeResponse(ctx context.Context, r interface{}) (response interface{}, err error) {
	return r, nil
}

// The service.Service method of AuthClientEndpoints is used for the endpoint of the client
type AuthClientEndpoints struct {
	LoginClientEndpoint       endpoint.Endpoint
	VerifyClientEndpoint      endpoint.Endpoint
	HealthCheckClientEndpoint endpoint.Endpoint
}

func NewAuthClientEndpoints(loginClientEP, verifyClientEP, HealthCheckClientEP endpoint.Endpoint) AuthClientEndpoints {
	return AuthClientEndpoints{
		LoginClientEndpoint:       loginClientEP,
		VerifyClientEndpoint:      verifyClientEP,
		HealthCheckClientEndpoint: HealthCheckClientEP,
	}
}

func (ue AuthClientEndpoints) Login(ctx context.Context, username, password string, unexpired bool) (string, error) {
	resp, err := ue.LoginClientEndpoint(ctx, &authpb.AuthRequest{
		Username:  username,
		Password:  password,
		Unexpired: unexpired,
	})
	if err != nil {
		return "", err
	}
	if response, ok := resp.(*authpb.AuthResponse); ok {
		return response.Token, nil
	} else {
		return "", ErrResponseNull
	}
}

func (ue AuthClientEndpoints) Verify(ctx context.Context, token string) (bool, error) {
	resp, err := ue.VerifyClientEndpoint(ctx, &authpb.VerifyRequest{Token: token})
	if err != nil {
		return false, err
	}
	if response, ok := resp.(*authpb.VerifyResponse); ok {
		return response.Passed, nil
	} else {
		return false, ErrResponseNull
	}
}

func makeLoginClientEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(conn,
		"authpb.AuthService",
		"Login",
		decodeRequest,
		encodeResponse,
		authpb.AuthResponse{}, // 写成request会导致panic，is *authpb.AuthResponse, not *authpb.AuthRequest
	).Endpoint()
}

func makeVerifyClientEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(conn,
		"authpb.AuthService",
		"Verify",
		decodeRequest,
		encodeResponse,
		authpb.VerifyResponse{}, // 写成request会导致panic，is *authpb.VerifyResponse, not *authpb.VerifyRequest
	).Endpoint()
}
