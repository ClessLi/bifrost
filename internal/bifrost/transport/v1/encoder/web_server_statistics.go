package encoder

import (
	"context"
	"encoding/json"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type webServerStatistics struct{}

var _ Encoder = webServerStatistics{}

func (e webServerStatistics) EncodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.Statistics:
		jdata, err := json.Marshal(r)
		if err != nil {
			return nil, code.WrapC(err, code.ErrEncodingFailed, "failed to encode response")
		}

		return &pbv1.Statistics{
			JsonData: jdata,
		}, nil
	default:
		return nil, code.WithCode(code.ErrEncodingFailed, "invalid web server statistics response: %v", r)
	}
}

func NewWebServerStatisticsEncoder() Encoder {
	return new(webServerStatistics)
}
