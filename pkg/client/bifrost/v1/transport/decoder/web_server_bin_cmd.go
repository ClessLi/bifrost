package decoder

import (
	"context"
	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerBinCMD struct{}

func (w webServerBinCMD) DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	switch resp := resp.(type) {
	case *pbv1.ExecuteResponse: // encode `Exec` response
		return &v1.ExecuteResponse{
			Successful:     resp.Successful,
			StandardOutput: resp.Stdout,
			StandardError:  resp.Stderr,
		}, nil
	default:
		return nil, errors.Errorf("invalid web server binary command response: %v", resp)
	}
}

var _ Decoder = webServerBinCMD{}
