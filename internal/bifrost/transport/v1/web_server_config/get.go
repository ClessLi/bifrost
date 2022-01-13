package web_server_config

import (
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
)

func (w *webServerConfigServer) GetServerNames(ctx context.Context, null *pbv1.Null) (*pbv1.ServerNames, error) {
	_, resp, err := w.handler.HandlerGetServerNames().ServeGRPC(ctx, null)
	if err != nil {
		return nil, err
	}
	response := resp.(*pbv1.ServerNames)
	return response, nil
}

func (w *webServerConfigServer) Get(r *pbv1.ServerName, stream pbv1.WebServerConfig_GetServer) error {

	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.ServerConfig)
	n := len(response.GetJsonData())
	for i := 0; i < n; i += w.options.ChunkSize {
		if n <= i+w.options.ChunkSize {
			err = stream.Send(&pbv1.ServerConfig{
				ServerName: response.GetServerName(),
				JsonData:   response.JsonData[i:],
			})
		} else {
			err = stream.Send(&pbv1.ServerConfig{
				ServerName: response.GetServerName(),
				JsonData:   response.JsonData[i : i+w.options.ChunkSize],
			})
		}
		if err != nil {
			return err
		}
	}

	return nil
}
