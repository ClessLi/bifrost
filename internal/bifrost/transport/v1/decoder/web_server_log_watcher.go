package decoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/yongPhone/bifrost/internal/pkg/code"
)

type webServerLogWatcher struct{}

var _ Decoder = webServerLogWatcher{}

func (w webServerLogWatcher) DecodeRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.LogWatchRequest: // decode `Watch` request
		return &v1.WebServerLogWatchRequest{
			ServerName:          &v1.ServerName{Name: r.ServerName},
			LogName:             r.LogName,
			FilteringRegexpRule: r.FilterRule,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerLogWatcherDecoder() Decoder {
	return new(webServerLogWatcher)
}
