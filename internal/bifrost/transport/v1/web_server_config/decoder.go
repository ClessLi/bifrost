package web_server_config

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	txpv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type decoder struct{}

var _ txpv1.Decoder = decoder{}

func (d decoder) DecodeRequest(_ context.Context, r interface{}) (interface{}, error) {
	switch r.(type) {
	case *pbv1.ServerName:
		return &v1.ServerName{Name: r.(*pbv1.ServerName).GetName()}, nil
	case *pbv1.ServerConfig:
		req := r.(*pbv1.ServerConfig)

		//srvconf, err := configuration.NewConfigurationFromJsonBytes(req.GetJsonData())
		//if err != nil {
		//	return nil, errors.WrapC(err, code.ErrDecodingFailed, "create configuration failed")
		//}

		return &v1.WebServerConfig{
			ServerName: &v1.ServerName{Name: req.GetServerName()},
			JsonData:   req.GetJsonData(),
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewDecoder() txpv1.Decoder {
	return &decoder{}
}
