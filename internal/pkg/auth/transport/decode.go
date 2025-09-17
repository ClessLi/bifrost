package transport

import (
	"errors"

	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/internal/pkg/auth/endpoint"

	"golang.org/x/net/context"
)

func DecodeAuthRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*authpb.AuthRequest); ok {
		return endpoint.AuthRequest{
			RequestType: "Login",
			Username:    req.Username,
			Password:    req.Password,
			Unexpired:   req.Unexpired,
		}, nil
	}
	return nil, errors.New("request has only one type: VerifyRequest")
}

func DecodeVerifyRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*authpb.VerifyRequest); ok {
		return endpoint.VerifyRequest{
			ResquesType: "Verify",
			Token:       req.Token,
		}, nil
	}
	return nil, errors.New("request has only one type: VerifyRequest")
}
