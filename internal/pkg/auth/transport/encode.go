package transport

import (
	"golang.org/x/net/context"

	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/internal/pkg/auth/endpoint"
)

func EncodeVerifyResponse(_ context.Context, r interface{}) (response interface{}, err error) {
	resp := r.(endpoint.VerifyResponse)
	if resp.Error != nil {
		//return &authpb.VerifyResponse{
		//	Passed: resp.Result,
		//	Err:    resp.Error.Error(),
		//}, resp.Error
		return nil, resp.Error
	}

	return &authpb.VerifyResponse{
		Passed: resp.Result,
		Err:    "",
	}, nil
}

func EncodeAuthResponse(_ context.Context, r interface{}) (response interface{}, err error) {
	resp := r.(endpoint.AuthResponse)
	if resp.Error != nil {
		//return &authpb.AuthResponse{
		//	Token: resp.Result,
		//	Err:   resp.Error.Error(),
		//}, resp.Error
		return nil, resp.Error
	}

	return &authpb.AuthResponse{
		Token: resp.Result,
		Err:   "",
	}, nil
}
