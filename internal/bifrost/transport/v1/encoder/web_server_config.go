package encoder

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
)

type webServerConfig struct{}

var _ Encoder = webServerConfig{}

func (e webServerConfig) EncodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.ServerNames: // encode `GetServerNames` response
		encodeServerNames := &pbv1.ServerNames{Names: make([]*pbv1.ServerName, 0)}
		for _, serverName := range *r {
			encodeServerNames.Names = append(encodeServerNames.Names, &pbv1.ServerName{Name: serverName.Name})
		}

		return encodeServerNames, nil
	case *v1.WebServerConfig: // encode `Get` response
		return &pbv1.ServerConfig{
			ServerName:           r.ServerName.Name,
			JsonData:             r.JsonData,
			OriginalFingerprints: r.OriginalFingerprints,
		}, nil
	case *v1.ContextData: // encode `ConnectivityCheckOfProxiedServers` response
		return &pbv1.ContextData{JsonData: r.JsonData}, nil
	case *v1.Response: // encode `Update` response
		return &pbv1.Response{Msg: []byte(r.Message)}, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server config response: %v", r)
	}
}

func NewWebServerConfigEncoder() Encoder {
	return new(webServerConfig)
}
