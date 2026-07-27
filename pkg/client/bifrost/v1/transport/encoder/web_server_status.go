package encoder

import (
	"context"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerStatus struct{}

func (w webServerStatus) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *pbv1.Null:
		return req, nil
	default:
		return nil, code.Errorf("invalid web server status request: %v", req)
	}
}

var _ Encoder = webServerStatus{}
