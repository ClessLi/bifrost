package encoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerStatistics struct{}

func (w webServerStatistics) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *v1.ServerName: // encode `Get` request
		return &pbv1.ServerName{Name: req.Name}, nil
	default:
		return nil, errors.Errorf("invalid web server statistics request: %v", req)
	}
}

var _ Encoder = webServerStatistics{}
