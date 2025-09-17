package auth

import (
	"errors"

	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var ErrResponseNull = errors.New("response is null")

func decodeRequest(_ context.Context, r interface{}) (request interface{}, err error) {
	return r, nil
}

func encodeResponse(_ context.Context, r interface{}) (response interface{}, err error) {
	return r, nil
}

// The service.Service method of authClientEndpoints is used for the endpoint of the client.
type authClientEndpoints struct {
	loginClientEndpoint  endpoint.Endpoint
	verifyClientEndpoint endpoint.Endpoint
}

func newAuthClientEndpoints(loginClientEP, verifyClientEP endpoint.Endpoint) *authClientEndpoints {
	return &authClientEndpoints{
		loginClientEndpoint:  loginClientEP,
		verifyClientEndpoint: verifyClientEP,
	}
}

func (ue authClientEndpoints) Login(ctx context.Context, username, password string, unexpired bool) (string, error) {
	resp, err := ue.loginClientEndpoint(ctx, &authpb.AuthRequest{
		Username:  username,
		Password:  password,
		Unexpired: unexpired,
	})
	if err != nil {
		return "", err
	}

	if response, ok := resp.(*authpb.AuthResponse); ok {
		return response.Token, nil
	}

	return "", ErrResponseNull
}

func (ue authClientEndpoints) Verify(ctx context.Context, token string) (bool, error) {
	resp, err := ue.verifyClientEndpoint(ctx, &authpb.VerifyRequest{Token: token})
	if err != nil {
		return false, err
	}

	if response, ok := resp.(*authpb.VerifyResponse); ok {
		return response.Passed, nil
	}

	return false, ErrResponseNull
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
