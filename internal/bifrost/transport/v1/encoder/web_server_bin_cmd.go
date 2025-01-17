package encoder

import (
	"context"
	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerBinCMD struct{}

var _ Encoder = &webServerBinCMD{}

func (w webServerBinCMD) EncodeResponse(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.ExecuteResponse:
		return &pbv1.ExecuteResponse{
			Successful: r.Successful,
			Stdout:     r.StandardOutput,
			Stderr:     r.StandardError,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server binary command executed response: %v", r)
	}
}

func NewWebServerBinCMDEncoder() Encoder {
	return new(webServerBinCMD)
}
