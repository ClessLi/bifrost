package encoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerConfig struct{}

func (w webServerConfig) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case nil: // encode `GetServerNames` request
		return &pbv1.Null{}, nil
	case *v1.ServerName: // encode `Get` request
		return &pbv1.ServerName{Name: req.Name}, nil
	case *v1.WebServerConfig: // encode `Update` request
		return &pbv1.ServerConfig{
			ServerName: req.ServerName.Name,
			JsonData:   req.JsonData,
		}, nil
	default:
		return nil, errors.Errorf("invalid web server config request: %v", req)
	}
}

var _ Encoder = webServerConfig{}
