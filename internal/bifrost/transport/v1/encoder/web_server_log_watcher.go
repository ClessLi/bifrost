package encoder

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type webServerLogWatcher struct{}

var _ Encoder = webServerLogWatcher{}

func (w webServerLogWatcher) EncodeResponse(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.WebServerLog: // return a bytes channel structure(point) *v1.WebServerLog from Watch endpoint, not a *v1.Response
		return r, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server log watcher response: %v", r)
	}
}

func NewWebServerLogWatcherEncoder() Encoder {
	return new(webServerLogWatcher)
}
