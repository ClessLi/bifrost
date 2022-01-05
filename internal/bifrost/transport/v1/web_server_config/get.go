package web_server_config

import pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

func (w *webServerConfigServer) Get(r *pbv1.ServerName, stream pbv1.WebServerConfig_GetServer) error {

	_, resp, err := w.handler.WebServerConfig().HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.ServerConfig)
	n := len(response.GetJsonData())
	for i := 0; i < n; i += w.Options.ChunkSize {
		if n <= i+w.Options.ChunkSize {
			err = stream.Send(&pbv1.ServerConfig{
				ServerName: response.GetServerName(),
				JsonData:   response.JsonData[i:],
			})
		} else {
			err = stream.Send(&pbv1.ServerConfig{
				ServerName: response.GetServerName(),
				JsonData:   response.JsonData[i : i+w.Options.ChunkSize],
			})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
