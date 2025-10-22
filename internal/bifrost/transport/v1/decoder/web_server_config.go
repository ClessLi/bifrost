package decoder

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
)

type webServerConfig struct{}

var _ Decoder = webServerConfig{}

func (d webServerConfig) DecodeRequest(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.Null: // decode `GetServerNames` request
		return r, nil
	case *pbv1.ServerName: // decode `Get` request
		return &v1.ServerName{Name: r.GetName()}, nil
	case *pbv1.ServerConfigContextPos: // decode `ConnectivityCheckOfProxiedServers` request
		return &v1.WebServerConfigContextPos{
			ServerName: &v1.ServerName{Name: r.ServerName},
			ContextPos: &v1.ContextPos{
				ConfigPath: r.ContextPos.ConfigPath,
				PosIndex:   r.ContextPos.Pos,
			},
			OriginalFingerprints: r.OriginalFingerprints,
		}, nil
	case *pbv1.ServerConfig: // decode `Update` request
		return &v1.WebServerConfig{
			ServerName:           &v1.ServerName{Name: r.GetServerName()},
			JsonData:             r.GetJsonData(),
			OriginalFingerprints: r.GetOriginalFingerprints(),
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerConfigDecoder() Decoder {
	return new(webServerConfig)
}
