package encoder

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerLogWatcher struct{}

func (w webServerLogWatcher) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *v1.WebServerLogWatchRequest: // encode `Watch` request
		return &pbv1.LogWatchRequest{
			ServerName: req.ServerName.Name,
			LogName:    req.LogName,
			FilterRule: req.FilteringRegexpRule,
		}, nil
	default:
		return nil, code.Errorf("invalid web server log watcher request: %v", req)
	}
}

var _ Encoder = webServerLogWatcher{}
