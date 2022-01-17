package decoder

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type webServerStatistics struct{}

var _ Decoder = webServerStatistics{}

func (d webServerStatistics) DecodeRequest(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.ServerName:
		return &v1.ServerName{Name: r.GetName()}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerStatisticsDecoder() Decoder {
	return new(webServerStatistics)
}
