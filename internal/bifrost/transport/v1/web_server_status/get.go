package web_server_status

import pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

func (w *webServerStatusServer) Get(r *pbv1.Null, stream pbv1.WebServerStatus_GetServer) error {

	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.Metrics)
	n := len(response.GetJsonData())
	for i := 0; i < n; i += w.options.ChunkSize - 2 {
		if n <= i+w.options.ChunkSize-2 {
			err = stream.Send(&pbv1.Metrics{JsonData: response.JsonData[i:]})
		} else {
			err = stream.Send(&pbv1.Metrics{JsonData: response.JsonData[i : i+w.options.ChunkSize-2]})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
