package encoder

import (
	"context"
	"encoding/json"

	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerStatus struct{}

var _ Encoder = &webServerStatus{}

func (w webServerStatus) EncodeResponse(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.Metrics:
		jdata, err := json.Marshal(r)
		if err != nil {
			return nil, errors.WithCode(code.ErrEncodingFailed, err.Error())
		}

		return &pbv1.Metrics{
			JsonData: jdata,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server status response: %v", r)
	}
}

func NewWebServerStatusEncoder() Encoder {
	return new(webServerStatus)
}
