package decoder

import (
	"context"
	"encoding/json"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/marmotedu/errors"
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
