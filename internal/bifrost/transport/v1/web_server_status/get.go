package web_server_status

import (
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/yongPhone/bifrost/internal/bifrost/transport/v1/utils"
)

func (w *webServerStatusServer) Get(r *pbv1.Null, stream pbv1.WebServerStatus_GetServer) error {
	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.Metrics)

	return utils.StreamSendMsg(stream, response.GetJsonData(), w.options.ChunkSize, func(msg []byte) interface{} {
		return &pbv1.Metrics{JsonData: msg}
	})
}
