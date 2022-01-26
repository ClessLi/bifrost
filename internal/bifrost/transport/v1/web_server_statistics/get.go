package web_server_statistics

import pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

func (w *webServerStatisticsServer) Get(r *pbv1.ServerName, stream pbv1.WebServerStatistics_GetServer) error {

	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.Statistics)
	n := len(response.GetJsonData())
	for i := 0; i < n; i += w.options.ChunkSize - 2 {
		if n <= i+w.options.ChunkSize-2 {
			err = stream.Send(&pbv1.Statistics{JsonData: response.JsonData[i:]})
		} else {
			err = stream.Send(&pbv1.Statistics{JsonData: response.JsonData[i : i+w.options.ChunkSize-2]})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
