package web_server_statistics

import (
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/yongPhone/bifrost/internal/bifrost/transport/v1/utils"
)

func (w *webServerStatisticsServer) Get(r *pbv1.ServerName, stream pbv1.WebServerStatistics_GetServer) error {
	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.Statistics)

	return utils.StreamSendMsg(stream, response.GetJsonData(), w.options.ChunkSize, func(msg []byte) interface{} {
		return &pbv1.Statistics{JsonData: msg}
	})
}
