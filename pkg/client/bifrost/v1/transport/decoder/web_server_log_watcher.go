package decoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

type webServerLogWatcher struct{}

func (w webServerLogWatcher) DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	switch resp := resp.(type) {
	case *v1.WebServerLog: // return a bytes channel structure(point) *v1.WebServerLog from Watch endpoint, not a *pbv1.Response
		return resp, nil
	default:
		return nil, errors.Errorf("invalid web server log watcher response: %v", resp)
	}
}

var _ Decoder = webServerLogWatcher{}
