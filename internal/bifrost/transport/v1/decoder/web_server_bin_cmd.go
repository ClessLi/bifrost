package decoder

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"

	"github.com/marmotedu/errors"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerBinCMD struct{}

var _ Decoder = webServerBinCMD{}

func (w webServerBinCMD) DecodeRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.ExecuteRequest:
		return &v1.ExecuteRequest{
			ServerName: r.ServerName,
			Args:       r.Args,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid execute request: %v", r)
	}
}

func NewWebServerBinCMDDecoder() Decoder {
	return new(webServerBinCMD)
}
