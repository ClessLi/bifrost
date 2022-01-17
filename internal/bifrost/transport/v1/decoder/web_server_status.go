package decoder

import (
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type webServerStatus struct{}

var _ Decoder = webServerStatus{}

func (w webServerStatus) DecodeRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.Null:
		return r, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerStatusDecoder() Decoder {
	return new(webServerStatus)
}
