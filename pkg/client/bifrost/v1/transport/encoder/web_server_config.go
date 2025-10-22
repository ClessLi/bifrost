package encoder

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	"github.com/marmotedu/errors"
)

type webServerConfig struct{}

func (w webServerConfig) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case nil: // encode `GetServerNames` request
		return &pbv1.Null{}, nil
	case *v1.ServerName: // encode `Get` request
		return &pbv1.ServerName{Name: req.Name}, nil
	case *v1.WebServerConfigContextPos: // encode `ConnectivityCheckOfProxiedServers` request
		return &pbv1.ServerConfigContextPos{
			ServerName: req.ServerName.Name,
			ContextPos: &pbv1.ContextPos{
				ConfigPath: req.ContextPos.ConfigPath,
				Pos:        req.ContextPos.PosIndex,
			},
			OriginalFingerprints: req.OriginalFingerprints,
		}, nil
	case *v1.WebServerConfig: // encode `Update` request
		return &pbv1.ServerConfig{
			ServerName:           req.ServerName.Name,
			JsonData:             req.JsonData,
			OriginalFingerprints: req.OriginalFingerprints,
		}, nil
	default:
		return nil, errors.Errorf("invalid web server config request: %v", req)
	}
}

var _ Encoder = webServerConfig{}
