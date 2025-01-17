package encoder

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	"github.com/marmotedu/errors"
)

type webServerBinCMD struct{}

func (w webServerBinCMD) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *v1.ExecuteRequest: // encode `Exec` request
		return &pbv1.ExecuteRequest{
			ServerName: req.ServerName,
			Args:       req.Args,
		}, nil
	default:
		return nil, errors.Errorf("invalid web server binary command request: %v", req)
	}
}

var _ Encoder = webServerBinCMD{}
