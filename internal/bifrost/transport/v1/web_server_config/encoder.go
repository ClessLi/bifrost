package web_server_config

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	txpv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type encoder struct{}

var _ txpv1.Encoder = encoder{}

func (e encoder) EncodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r.(type) {
	case *v1.WebServerConfig:
		resp := r.(*v1.WebServerConfig)
		return &pbv1.ServerConfig{
			ServerName: resp.ServerName.Name,
			JsonData:   resp.JsonData,
		}, nil
	case *v1.Response:
		resp := r.(*v1.Response)
		return &pbv1.Response{Msg: []byte(resp.Message)}, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server config response: %v", r)
	}
}

func NewEncoder() txpv1.Encoder {
	return &encoder{}
}
