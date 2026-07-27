package encoder

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerStatistics struct{}

func (w webServerStatistics) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *v1.ServerName: // encode `Get` request
		return &pbv1.ServerName{Name: req.Name}, nil
	default:
		return nil, code.Errorf("invalid web server statistics request: %v", req)
	}
}

var _ Encoder = webServerStatistics{}
