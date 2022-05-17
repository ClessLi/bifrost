package encoder

import (
	"context"
	"encoding/json"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/yongPhone/bifrost/internal/pkg/code"
)

type webServerStatistics struct{}

var _ Encoder = webServerStatistics{}

func (e webServerStatistics) EncodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *v1.Statistics:
		jdata, err := json.Marshal(r)
		if err != nil {
			return nil, errors.WithCode(code.ErrEncodingFailed, err.Error())
		}

		return &pbv1.Statistics{
			JsonData: jdata,
		}, nil
	default:
		return nil, errors.WithCode(code.ErrEncodingFailed, "invalid web server statistics response: %v", r)
	}
}

func NewWebServerStatisticsEncoder() Encoder {
	return new(webServerStatistics)
}
