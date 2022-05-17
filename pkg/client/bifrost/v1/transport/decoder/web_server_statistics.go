package decoder

import (
	"context"
	"encoding/json"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerStatistics struct{}

func (w webServerStatistics) DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	switch resp := resp.(type) {
	case *pbv1.Statistics:
		statistics := new(v1.Statistics)
		err := json.Unmarshal(resp.GetJsonData(), statistics)

		return statistics, err
	default:
		return nil, errors.Errorf("invalid web server statistics response: %v", resp)
	}
}

var _ Decoder = webServerStatistics{}
