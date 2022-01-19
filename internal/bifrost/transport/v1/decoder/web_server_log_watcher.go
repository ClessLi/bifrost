package decoder

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type webServerLogWatcher struct{}

var _ Decoder = webServerLogWatcher{}

func (w webServerLogWatcher) DecodeRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.LogWatchRequest:
		return &v1.WebServerLogWatchRequest{
			ServerName:          &v1.ServerName{Name: r.ServerName},
			LogPath:             r.LogName,
			FilteringRegexpRule: r.FilterRule,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerLogWatcherDecoder() Decoder {
	return new(webServerLogWatcher)
}
